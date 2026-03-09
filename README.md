<p align="center">
  <img src="https://picui.ogmua.cn/s1/2026/03/09/69ae6edbd4635.webp" alt="FRP Desktop Logo" width="100">
</p>

<h1 align="center">FRP Desktop</h1>

<p align="center">
  <strong>一款面向普通用户的 FRP 图形客户端</strong>
</p>

<p align="center">
  <a href="https://github.com/dsfeiji/frp-desktop/releases">
    <img src="https://img.shields.io/github/v/release/dsfeiji/frp-desktop?style=flat-square&color=blue" alt="Release">
  </a>
  <img src="https://img.shields.io/github/license/dsfeiji/frp-desktop?style=flat-square" alt="License">
  <img src="https://img.shields.io/github/stars/dsfeiji/frp-desktop?style=flat-square" alt="Stars">
  <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20macOS-brightgreen?style=flat-square" alt="Platform">
</p>

<p align="center">
  <img src="https://picui.ogmua.cn/s1/2026/03/09/69ae6f91204c4.webp" alt="FRP Desktop 软件界面" width="700" style="border-radius: 8px;">
</p>

---

FRP Desktop 让你告别手写 `frpc.toml`。只需填写服务器信息和本地端口，即可一键启动端口转发，实现快速的内网穿透。

## 🛠️ 这个软件能做什么

- **图形化配置**：轻松填写 FRP 服务器（地址、端口、Token）。
- **一键启停**：快速启动或停止端口转发任务。
- **多端口映射**：支持同时转发多个本地端口（如 `3000, 8080`）。
- **自动托管**：自动生成并管理底层 `frpc` 配置文件。
- **实时监控**：显示运行状态、进程 PID 及错误反馈。

## 🌐 适用场景

- 你有一台已部署 `frps` 的服务器。
- 你希望把本机服务（如 3000、8080、5000 端口）映射到公网。
- 你不想每次都手写复杂的命令和配置文件。

## 📦 安装指南

请在 [GitHub Releases](https://github.com/dsfeiji/frp-desktop/releases) 页面下载对应版本：

- **macOS**: 下载 `.dmg` 文件安装。
- **Windows**: 下载 `.exe` 安装包或 `.zip` 绿色版。

### 🔗 下载链接（直达）

- Release 页面：<https://github.com/dsfeiji/frp-desktop/releases>
- 最新版本：<https://github.com/dsfeiji/frp-desktop/releases/latest>

## 💡 使用方法（只需 3 步）

1. **配置服务器**：点击右上角齿轮，填写地址、端口（通常 7000）和 Token。
2. **输入端口**：在主界面输入需要转发的本地端口，例如：`3000, 8080`。
3. **开启转发**：点击 **[开始转发]**，状态显示“运行中”即成功。

## 🧩 服务器脚本（自动安装 frps）

如果你还没有部署 `frps`，可以直接使用仓库内脚本一键安装。  
脚本会自动下载 `frps`、生成 `frps.toml`、随机生成 **8 位字母数字 Token**、并启动服务。

### Linux 一键执行

脚本下载链接：  
<https://raw.githubusercontent.com/dsfeiji/frp-desktop/main/scripts/gen_frps_config.sh>

执行命令（root）：

```bash
curl -fsSL https://raw.githubusercontent.com/dsfeiji/frp-desktop/main/scripts/gen_frps_config.sh -o gen_frps_config.sh
sudo sh gen_frps_config.sh
```

### Windows 一键执行

脚本下载链接：  
<https://raw.githubusercontent.com/dsfeiji/frp-desktop/main/scripts/gen_frps_config.bat>

执行方式（管理员 CMD）：

```bat
powershell -Command "Invoke-WebRequest -Uri 'https://raw.githubusercontent.com/dsfeiji/frp-desktop/main/scripts/gen_frps_config.bat' -OutFile 'gen_frps_config.bat'"
gen_frps_config.bat
```

## ❓ 常见问题

- **1) 启动失败？** 请检查服务器配置（地址/端口/Token）是否正确，以及服务器端 `frps` 是否在线。
- **2) macOS 无法打开？** 因未做开发者签名，请通过“右键图标 -> 打开”或在“系统设置 -> 隐私与安全性”中选择“仍要打开”。
- **3) 配置会保存吗？** 会。软件会自动保存你的所有设置，下次打开可直接使用。
- **4) 忘记 Token 去哪里看？**  
  - Linux（默认脚本安装）：`/etc/frp/frps.toml`，可执行 `cat /etc/frp/frps.toml` 查看 `token`。  
  - Windows（默认脚本安装）：`C:\frp-server\frps.toml`，可用记事本打开查看 `token`。  
  - 也可以重跑服务器脚本生成新 Token，然后在客户端同步更新。

---

> [!IMPORTANT]
> **安全提示**：请勿在公开仓库或截图中泄露真实服务器地址与 Token。

**项目仓库:** [https://github.com/dsfeiji/frp-desktop](https://github.com/dsfeiji/frp-desktop)
