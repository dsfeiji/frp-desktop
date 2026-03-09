import './style.css';

import { EventsOn } from '../wailsjs/runtime/runtime';
import {
  GetConfig,
  GetRuntimeState,
  SaveConfig,
  StartFrpc,
  StopFrpc,
} from '../wailsjs/go/main/App';

const appRoot = document.querySelector('#app');

appRoot.innerHTML = `
  <div class="shell">
    <main class="app-frame">
      <header class="topbar">
        <div class="brand">
          <p class="eyebrow">FRP Desktop</p>
          <h1>端口转发</h1>
        </div>
        <button id="settingsBtn" class="icon-btn" title="服务器设置" aria-label="服务器设置">
          <span aria-hidden="true">⚙</span>
        </button>
      </header>

      <section class="panel panel-main">
        <div class="panel-head">
          <h2>本地端口</h2>
          <span class="pill">TCP</span>
        </div>

        <label for="portsInput">端口列表</label>
        <input id="portsInput" type="text" placeholder="例如: 80,443,8080" />
        <p class="muted">范围 1-65535，多个端口请用英文逗号分隔。</p>
      </section>

      <section class="panel panel-actions">
        <button id="startBtn" class="btn btn-primary">开始转发</button>
        <button id="stopBtn" class="btn btn-secondary">停止</button>
        <p id="status" class="status">状态: 未运行</p>
      </section>
    </main>
  </div>

  <div id="settingsModal" class="modal hidden" role="dialog" aria-modal="true" aria-labelledby="settingsTitle">
    <div class="modal-card">
      <div class="panel-head">
        <h2 id="settingsTitle">服务器设置</h2>
      </div>

      <label for="serverAddrInput">服务器地址</label>
      <input id="serverAddrInput" type="text" placeholder="例如: frp.example.com" />

      <label for="serverPortInput">服务器端口</label>
      <input id="serverPortInput" type="number" min="1" max="65535" placeholder="例如: 7000" />

      <label for="authTokenInput">认证 Token</label>
      <input id="authTokenInput" type="text" placeholder="请输入 token" />

      <div class="modal-actions">
        <button id="saveSettingsBtn" class="btn btn-primary">保存</button>
        <button id="cancelSettingsBtn" class="btn btn-secondary">取消</button>
      </div>
    </div>
  </div>
`;

const currentConfig = {
  frpcPath: '',
  localPorts: [],
  serverAddr: '',
  serverPort: 7000,
  authToken: '',
};

function setStatus(text, isErr = false) {
  const el = document.getElementById('status');
  el.textContent = text;
  el.classList.toggle('err', isErr);
}

function openSettings() {
  const modal = document.getElementById('settingsModal');
  modal.classList.remove('hidden');
}

function closeSettings() {
  const modal = document.getElementById('settingsModal');
  modal.classList.add('hidden');
}

function syncSettingsInputs() {
  document.getElementById('serverAddrInput').value = currentConfig.serverAddr || '';
  document.getElementById('serverPortInput').value = currentConfig.serverPort || 7000;
  document.getElementById('authTokenInput').value = currentConfig.authToken || '';
}

async function refreshState() {
  const state = await GetRuntimeState();
  if (state.running) {
    setStatus(`状态: 运行中 (PID ${state.pid})`);
  } else if (state.lastError) {
    setStatus(`状态: 未运行 (${state.lastError})`, true);
  } else {
    setStatus('状态: 未运行');
  }
}

async function saveServerSettings() {
  const serverAddr = document.getElementById('serverAddrInput').value.trim();
  const serverPort = Number.parseInt(document.getElementById('serverPortInput').value.trim(), 10);
  const authToken = document.getElementById('authTokenInput').value.trim();

  if (!serverAddr) {
    setStatus('状态: 服务器地址不能为空', true);
    return;
  }
  if (!Number.isInteger(serverPort) || serverPort <= 0 || serverPort > 65535) {
    setStatus('状态: 服务器端口无效 (1-65535)', true);
    return;
  }

  currentConfig.serverAddr = serverAddr;
  currentConfig.serverPort = serverPort;
  currentConfig.authToken = authToken;

  await SaveConfig(currentConfig);
  closeSettings();
  setStatus('状态: 服务器设置已保存');
}

async function startForward() {
  const raw = document.getElementById('portsInput').value.trim();
  const ports = raw
    .split(',')
    .map((s) => Number.parseInt(s.trim(), 10))
    .filter((n) => Number.isInteger(n) && n > 0 && n <= 65535);

  const uniquePorts = [...new Set(ports)].sort((a, b) => a - b);
  if (!uniquePorts.length) {
    setStatus('状态: 请输入有效端口，示例 80,443,8080', true);
    return;
  }

  currentConfig.localPorts = uniquePorts;
  setStatus('状态: 启动中...');
  await SaveConfig(currentConfig);
  await StartFrpc();
  await refreshState();
}

async function boot() {
  const cfg = await GetConfig();
  Object.assign(currentConfig, cfg);

  if (!Array.isArray(currentConfig.localPorts)) {
    currentConfig.localPorts = [];
  }
  if (!currentConfig.serverAddr) {
    currentConfig.serverAddr = '';
  }
  if (!Number.isInteger(currentConfig.serverPort) || currentConfig.serverPort <= 0 || currentConfig.serverPort > 65535) {
    currentConfig.serverPort = 7000;
  }
  if (!currentConfig.authToken) {
    currentConfig.authToken = '';
  }

  if (currentConfig.localPorts.length) {
    document.getElementById('portsInput').value = currentConfig.localPorts.join(',');
  }
  syncSettingsInputs();
  await refreshState();

  document.getElementById('settingsBtn').addEventListener('click', () => {
    syncSettingsInputs();
    openSettings();
  });

  document.getElementById('cancelSettingsBtn').addEventListener('click', () => {
    closeSettings();
  });

  document.getElementById('saveSettingsBtn').addEventListener('click', async () => {
    try {
      await saveServerSettings();
    } catch (err) {
      setStatus(`状态: 保存设置失败 - ${String(err)}`, true);
    }
  });

  document.getElementById('settingsModal').addEventListener('click', (event) => {
    if (event.target.id === 'settingsModal') {
      closeSettings();
    }
  });

  document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') {
      closeSettings();
    }
  });

  document.getElementById('startBtn').addEventListener('click', async () => {
    try {
      await startForward();
    } catch (err) {
      setStatus(`状态: 启动失败 - ${String(err)}`, true);
    }
  });

  document.getElementById('stopBtn').addEventListener('click', async () => {
    try {
      await StopFrpc();
      setStatus('状态: 停止中...');
      setTimeout(() => {
        refreshState().catch(() => {});
      }, 300);
    } catch (err) {
      setStatus(`状态: 停止失败 - ${String(err)}`, true);
    }
  });

  EventsOn('frpc:status', (state) => {
    if (state.running) {
      setStatus(`状态: 运行中 (PID ${state.pid})`);
    } else if (state.lastError) {
      setStatus(`状态: 未运行 (${state.lastError})`, true);
    } else {
      setStatus('状态: 未运行');
    }
  });

  setInterval(() => {
    refreshState().catch(() => {});
  }, 1500);
}

boot().catch((err) => {
  setStatus(`状态: 初始化失败 - ${String(err)}`, true);
});
