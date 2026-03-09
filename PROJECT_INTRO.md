# FRP Desktop

FRP Desktop is a lightweight desktop client for FRP (`frpc`) built with **Go + Wails**.
It provides a clean UI to configure server connection details, select local ports, and control forwarding without editing config files manually.

## Why this project

- Reduce manual `frpc.toml` editing
- Provide a simple workflow for non-terminal users
- Offer cross-platform desktop usage with a consistent interface

## Core Capabilities

- Server settings editor (address, port, token)
- Multi-port forwarding start/stop controls
- Runtime status and error feedback
- Automatic `frpc.toml` generation
- Optional bundled `frpc` in packaged app builds

## Tech Stack

- Backend: Go + Wails v2
- Frontend: Vite + Vanilla JavaScript + CSS
- Packaging: Wails build pipeline (desktop app output)

## Typical Use Flow

1. Open app and configure FRP server settings.
2. Enter one or more local ports.
3. Click **Start Forwarding**.
4. Monitor status and logs in the app.

## Repository Notes

This repository is prepared for open-source hosting.  
All sensitive defaults should remain placeholders. Real deployment values must be configured locally.
