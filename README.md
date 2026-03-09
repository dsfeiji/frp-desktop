# FRP Desktop

[English](#english) | [中文](#中文)

## English

FRP Desktop is a lightweight desktop client for FRP (`frpc`) built with **Go + Wails**.

- Project introduction: [PROJECT_INTRO.md](./PROJECT_INTRO.md)
- Chinese introduction: [PROJECT_INTRO_CN.md](./PROJECT_INTRO_CN.md)

### Quick Start

```bash
cd /Users/dsfeiji/Downloads/frp
PATH="$HOME/go/bin:$PATH" GOPROXY="https://goproxy.cn,direct" wails dev
```

### Build

```bash
PATH="$HOME/go/bin:$PATH" GOPROXY="https://goproxy.cn,direct" wails build -clean -platform darwin/universal
```

## 中文

FRP Desktop 是一个基于 **Go + Wails** 的 FRP（`frpc`）桌面客户端。

- 中文介绍文档：[PROJECT_INTRO_CN.md](./PROJECT_INTRO_CN.md)
- English intro: [PROJECT_INTRO.md](./PROJECT_INTRO.md)

### 快速开始

```bash
cd /Users/dsfeiji/Downloads/frp
PATH="$HOME/go/bin:$PATH" GOPROXY="https://goproxy.cn,direct" wails dev
```

### 构建

```bash
PATH="$HOME/go/bin:$PATH" GOPROXY="https://goproxy.cn,direct" wails build -clean -platform darwin/universal
```

## Security Note / 安全说明

Do not commit real server addresses, tokens, or private credentials.
请勿提交真实服务器地址、Token 或私密凭据到仓库。
