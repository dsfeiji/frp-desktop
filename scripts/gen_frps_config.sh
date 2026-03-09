#!/usr/bin/env sh
# One-click FRP server installer for Linux.
# It installs frps, writes /etc/frp/frps.toml, creates systemd service, and starts it.
#
# Usage:
#   sudo sh scripts/gen_frps_config.sh [frp_version] [bind_port] [install_dir] [config_dir]
#
# Example:
#   sudo sh scripts/gen_frps_config.sh 0.67.0 7000 /opt/frp /etc/frp

set -eu

FRP_VERSION="${1:-0.67.0}"
BIND_PORT="${2:-7000}"
INSTALL_DIR="${3:-/opt/frp}"
CONFIG_DIR="${4:-/etc/frp}"

if [ "$(id -u)" -ne 0 ]; then
  echo "Please run as root: sudo sh scripts/gen_frps_config.sh" >&2
  exit 1
fi

case "$BIND_PORT" in
  ''|*[!0-9]*) echo "bind_port must be numeric" >&2; exit 1 ;;
esac
if [ "$BIND_PORT" -lt 1 ] || [ "$BIND_PORT" -gt 65535 ]; then
  echo "bind_port must be 1-65535" >&2
  exit 1
fi

if command -v apt-get >/dev/null 2>&1; then
  apt-get update -y
  apt-get install -y curl tar
elif command -v yum >/dev/null 2>&1; then
  yum install -y curl tar
elif command -v dnf >/dev/null 2>&1; then
  dnf install -y curl tar
else
  echo "Please install curl and tar manually." >&2
fi

ARCH_RAW="$(uname -m)"
case "$ARCH_RAW" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH_RAW" >&2
    exit 1
    ;;
esac

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

PKG="frp_${FRP_VERSION}_linux_${ARCH}.tar.gz"
URL="https://github.com/fatedier/frp/releases/download/v${FRP_VERSION}/${PKG}"

echo "Downloading: $URL"
curl -fL "$URL" -o "$TMP_DIR/$PKG"
tar -xzf "$TMP_DIR/$PKG" -C "$TMP_DIR"

EXTRACT_DIR="$TMP_DIR/frp_${FRP_VERSION}_linux_${ARCH}"
if [ ! -f "$EXTRACT_DIR/frps" ]; then
  echo "frps binary not found in package" >&2
  exit 1
fi

mkdir -p "$INSTALL_DIR" "$CONFIG_DIR"
install -m 755 "$EXTRACT_DIR/frps" "$INSTALL_DIR/frps"

TOKEN="$(tr -dc 'A-Za-z0-9' </dev/urandom | head -c 8)"

cat >"$CONFIG_DIR/frps.toml" <<CONF
bindAddr = "0.0.0.0"
bindPort = ${BIND_PORT}

[auth]
method = "token"
token = "${TOKEN}"
CONF

cat >/etc/systemd/system/frps.service <<SERVICE
[Unit]
Description=FRP Server Service
After=network.target

[Service]
Type=simple
ExecStart=${INSTALL_DIR}/frps -c ${CONFIG_DIR}/frps.toml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
SERVICE

systemctl daemon-reload
systemctl enable frps
systemctl restart frps

echo "========================================"
echo "FRP server installed and started."
echo "frps binary : ${INSTALL_DIR}/frps"
echo "config file : ${CONFIG_DIR}/frps.toml"
echo "service     : frps"
echo "token       : ${TOKEN}"
echo "port        : ${BIND_PORT}"
echo "status      : systemctl status frps"
echo "========================================"
