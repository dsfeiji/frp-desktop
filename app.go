package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const maxLogLines = 500
const defaultServerAddr = ""
const defaultServerPort = 7000
const defaultAuthToken = ""

type AppConfig struct {
	FrpcPath   string `json:"frpcPath"`
	LocalPorts []int  `json:"localPorts"`
	LocalPort  int    `json:"localPort,omitempty"`
	ServerAddr string `json:"serverAddr"`
	ServerPort int    `json:"serverPort"`
	AuthToken  string `json:"authToken"`
}

type RuntimeState struct {
	Running    bool   `json:"running"`
	PID        int    `json:"pid"`
	AppDir     string `json:"appDir"`
	ConfigFile string `json:"configFile"`
	FrpcToml   string `json:"frpcToml"`
	FrpcPath   string `json:"frpcPath"`
	LastError  string `json:"lastError"`
}

type App struct {
	ctx context.Context

	mu         sync.Mutex
	config     AppConfig
	proc       *exec.Cmd
	logs       []string
	appDir     string
	configFile string
	frpcToml   string
	lastError  string
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	a.appDir = filepath.Join(dir, "frp-desktop")
	a.configFile = filepath.Join(a.appDir, "config.json")
	a.frpcToml = filepath.Join(a.appDir, "frpc.toml")

	if err := os.MkdirAll(a.appDir, 0o755); err != nil {
		a.appendLog("failed to create app dir: " + err.Error())
	}

	cfg, err := a.loadConfig()
	if err != nil {
		a.appendLog("failed to load config: " + err.Error())
		cfg = defaultConfig()
	}

	a.mu.Lock()
	a.config = normalizeConfig(cfg)
	a.mu.Unlock()

	if err := a.persistConfig(); err != nil {
		a.appendLog("failed to persist initial config: " + err.Error())
	}
}

func defaultConfig() AppConfig {
	return AppConfig{
		FrpcPath:   "",
		LocalPorts: []int{},
		LocalPort:  0,
		ServerAddr: defaultServerAddr,
		ServerPort: defaultServerPort,
		AuthToken:  defaultAuthToken,
	}
}

func normalizeConfig(cfg AppConfig) AppConfig {
	// backward compatibility for previous single-port config
	if len(cfg.LocalPorts) == 0 && cfg.LocalPort > 0 && cfg.LocalPort <= 65535 {
		cfg.LocalPorts = []int{cfg.LocalPort}
	}
	cfg.LocalPorts = uniqueSortedPorts(cfg.LocalPorts)
	cfg.LocalPort = 0
	cfg.FrpcPath = strings.TrimSpace(cfg.FrpcPath)
	cfg.ServerAddr = strings.TrimSpace(cfg.ServerAddr)

	if cfg.ServerPort <= 0 || cfg.ServerPort > 65535 {
		cfg.ServerPort = defaultServerPort
	}
	cfg.AuthToken = strings.TrimSpace(cfg.AuthToken)
	if cfg.AuthToken == "" {
		cfg.AuthToken = defaultAuthToken
	}
	return cfg
}

func (a *App) loadConfig() (AppConfig, error) {
	b, err := os.ReadFile(a.configFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultConfig(), nil
		}
		return AppConfig{}, err
	}
	var cfg AppConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		return AppConfig{}, err
	}
	return normalizeConfig(cfg), nil
}

func (a *App) persistConfig() error {
	a.mu.Lock()
	cfg := a.config
	a.mu.Unlock()

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(a.configFile, b, 0o644)
}

func (a *App) GetConfig() AppConfig {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.config
}

func (a *App) SaveConfig(cfg AppConfig) error {
	a.mu.Lock()
	a.config = normalizeConfig(cfg)
	a.mu.Unlock()
	return a.persistConfig()
}

func (a *App) GetRuntimeState() RuntimeState {
	a.mu.Lock()
	defer a.mu.Unlock()

	state := RuntimeState{
		Running:    a.proc != nil,
		AppDir:     a.appDir,
		ConfigFile: a.configFile,
		FrpcToml:   a.frpcToml,
		LastError:  a.lastError,
	}
	if a.proc != nil && a.proc.Process != nil {
		state.PID = a.proc.Process.Pid
	}
	if p, _ := a.resolveFrpcPathLocked(); p != "" {
		state.FrpcPath = p
	}
	return state
}

func (a *App) GetLogs() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	cp := make([]string, len(a.logs))
	copy(cp, a.logs)
	return cp
}

func (a *App) ClearLogs() {
	a.mu.Lock()
	a.logs = nil
	a.mu.Unlock()
}

func (a *App) DiscoverPorts() ([]int, error) {
	return discoverListeningPorts()
}

func (a *App) StartFrpc() error {
	a.mu.Lock()
	if a.proc != nil {
		a.mu.Unlock()
		return errors.New("frpc is already running")
	}

	cfg := a.config
	frpcPath, err := a.resolveFrpcPathLocked()
	if err != nil {
		a.lastError = err.Error()
		a.mu.Unlock()
		return err
	}

	if len(cfg.LocalPorts) == 0 {
		err := errors.New("please input local ports")
		a.lastError = err.Error()
		a.mu.Unlock()
		return err
	}

	if strings.TrimSpace(cfg.ServerAddr) == "" {
		err := errors.New("please input server address")
		a.lastError = err.Error()
		a.mu.Unlock()
		return err
	}

	if cfg.ServerPort <= 0 || cfg.ServerPort > 65535 {
		err := errors.New("please input valid server port")
		a.lastError = err.Error()
		a.mu.Unlock()
		return err
	}

	if err := a.writeFrpcToml(cfg); err != nil {
		a.lastError = err.Error()
		a.mu.Unlock()
		return err
	}

	cmd := exec.Command(frpcPath, "-c", a.frpcToml)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		a.lastError = err.Error()
		a.mu.Unlock()
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		a.lastError = err.Error()
		a.mu.Unlock()
		return err
	}

	if err := cmd.Start(); err != nil {
		a.lastError = err.Error()
		a.mu.Unlock()
		return err
	}

	a.proc = cmd
	a.lastError = ""
	a.mu.Unlock()

	a.appendLog(fmt.Sprintf("frpc started. pid=%d localPorts=%v", cmd.Process.Pid, cfg.LocalPorts))
	a.emitStatus()

	go a.capturePipe(stdout, "stdout")
	go a.capturePipe(stderr, "stderr")

	go func() {
		err := cmd.Wait()
		a.mu.Lock()
		a.proc = nil
		if err != nil {
			a.lastError = err.Error()
		}
		a.mu.Unlock()

		if err != nil {
			a.appendLog("frpc exited with error: " + err.Error())
		} else {
			a.appendLog("frpc exited")
		}
		a.emitStatus()
	}()

	return nil
}

func (a *App) StopFrpc() error {
	a.mu.Lock()
	proc := a.proc
	a.mu.Unlock()

	if proc == nil || proc.Process == nil {
		return nil
	}

	if err := proc.Process.Kill(); err != nil {
		return err
	}
	a.appendLog("stop signal sent")
	return nil
}

func (a *App) writeFrpcToml(cfg AppConfig) error {
	var b strings.Builder

	b.WriteString("serverAddr = \"")
	b.WriteString(escapeTomlString(cfg.ServerAddr))
	b.WriteString("\"\n")
	b.WriteString("serverPort = ")
	b.WriteString(strconv.Itoa(cfg.ServerPort))
	b.WriteString("\n")
	b.WriteString("transport.protocol = \"quic\"\n\n")

	if cfg.AuthToken != "" {
		b.WriteString("[auth]\n")
		b.WriteString("method = \"token\"\n")
		b.WriteString("token = \"")
		b.WriteString(escapeTomlString(cfg.AuthToken))
		b.WriteString("\"\n\n")
	}

	for _, port := range cfg.LocalPorts {
		b.WriteString("[[proxies]]\n")
		b.WriteString("name = \"tcp-")
		b.WriteString(strconv.Itoa(port))
		b.WriteString("\"\n")
		b.WriteString("type = \"tcp\"\n")
		b.WriteString("localIP = \"127.0.0.1\"\n")
		b.WriteString("localPort = ")
		b.WriteString(strconv.Itoa(port))
		b.WriteString("\n")
		b.WriteString("remotePort = ")
		b.WriteString(strconv.Itoa(port))
		b.WriteString("\n\n")
	}

	return os.WriteFile(a.frpcToml, []byte(b.String()), 0o644)
}

func (a *App) resolveFrpcPathLocked() (string, error) {
	cfg := a.config

	var candidates []string
	if cfg.FrpcPath != "" {
		candidates = append(candidates, cfg.FrpcPath)
	}

	exeName := "frpc"
	if runtime.GOOS == "windows" {
		exeName = "frpc.exe"
	}

	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, exeName),
			filepath.Join(exeDir, "frpc", exeName),
			filepath.Join(exeDir, "bin", exeName),
		)
	}

	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(wd, exeName),
			filepath.Join(wd, "frpc", exeName),
			filepath.Join(wd, "bin", exeName),
		)
	}

	for _, c := range candidates {
		if c == "" {
			continue
		}
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}

	if p, err := exec.LookPath(exeName); err == nil {
		return p, nil
	}

	return "", fmt.Errorf("frpc not found. set FrpcPath in config file or place %s in build/bin", exeName)
}

func discoverListeningPorts() ([]int, error) {
	cmd := exec.Command("netstat", "-an")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	set := map[int]struct{}{}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		upper := strings.ToUpper(line)
		if !strings.Contains(upper, "LISTEN") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		local := ""
		proto := strings.ToLower(fields[0])
		switch {
		case strings.HasPrefix(proto, "tcp"):
			if runtime.GOOS == "windows" {
				local = fields[1]
			} else {
				local = fields[3]
			}
		default:
			continue
		}

		if p, ok := parsePortFromAddress(local); ok {
			set[p] = struct{}{}
		}
	}

	ports := make([]int, 0, len(set))
	for p := range set {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	return ports, nil
}

func parsePortFromAddress(addr string) (int, bool) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return 0, false
	}

	var idx int
	if strings.Contains(addr, "]:") {
		idx = strings.LastIndex(addr, "]:")
		if idx >= 0 {
			addr = addr[idx+2:]
		}
	} else {
		idx = strings.LastIndex(addr, ":")
		if idx >= 0 {
			addr = addr[idx+1:]
		} else {
			idx = strings.LastIndex(addr, ".")
			if idx >= 0 {
				addr = addr[idx+1:]
			}
		}
	}

	p, err := strconv.Atoi(addr)
	if err != nil || p <= 0 || p > 65535 {
		return 0, false
	}
	return p, true
}

func uniqueSortedPorts(in []int) []int {
	set := map[int]struct{}{}
	for _, p := range in {
		if p <= 0 || p > 65535 {
			continue
		}
		set[p] = struct{}{}
	}
	out := make([]int, 0, len(set))
	for p := range set {
		out = append(out, p)
	}
	sort.Ints(out)
	return out
}

func escapeTomlString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

func (a *App) capturePipe(rdr interface{ Read([]byte) (int, error) }, stream string) {
	scanner := bufio.NewScanner(rdr)
	for scanner.Scan() {
		a.appendLog("[" + stream + "] " + scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		a.appendLog("[" + stream + "] read error: " + err.Error())
	}
}

func (a *App) appendLog(msg string) {
	line := time.Now().Format("15:04:05") + " " + msg

	a.mu.Lock()
	a.logs = append(a.logs, line)
	if len(a.logs) > maxLogLines {
		a.logs = append([]string(nil), a.logs[len(a.logs)-maxLogLines:]...)
	}
	a.mu.Unlock()

	if a.ctx != nil {
		wruntime.EventsEmit(a.ctx, "frpc:log", line)
	}
}

func (a *App) emitStatus() {
	if a.ctx == nil {
		return
	}
	wruntime.EventsEmit(a.ctx, "frpc:status", a.GetRuntimeState())
}
