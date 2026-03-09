#!/usr/bin/env sh
# Generate frps.toml for Linux/macOS.
#
# Usage:
#   sh scripts/gen_frps_config.sh [output] [server_addr] [bind_port] [token] [dashboard]
#                                 [dashboard_addr] [dashboard_port] [dashboard_user] [dashboard_password]
#
# Example:
#   sh scripts/gen_frps_config.sh ./frps.toml 0.0.0.0 7000 my_token 1 0.0.0.0 7500 admin admin123

set -eu

OUTPUT="${1:-frps.toml}"
SERVER_ADDR="${2:-0.0.0.0}"
BIND_PORT="${3:-7000}"
TOKEN="${4:-}"
DASHBOARD="${5:-0}"
DASHBOARD_ADDR="${6:-0.0.0.0}"
DASHBOARD_PORT="${7:-7500}"
DASHBOARD_USER="${8:-admin}"
DASHBOARD_PASSWORD="${9:-}"

is_valid_port() {
  case "$1" in
    ''|*[!0-9]*) return 1 ;;
  esac
  [ "$1" -ge 1 ] && [ "$1" -le 65535 ]
}

random_token() {
  # 32-char alphanumeric token.
  tr -dc 'A-Za-z0-9' </dev/urandom | head -c 32
}

if ! is_valid_port "$BIND_PORT"; then
  echo "bind_port must be 1-65535" >&2
  exit 1
fi

if [ "$DASHBOARD" = "1" ] && ! is_valid_port "$DASHBOARD_PORT"; then
  echo "dashboard_port must be 1-65535" >&2
  exit 1
fi

if [ -z "$TOKEN" ]; then
  TOKEN="$(random_token)"
fi

if [ "$DASHBOARD" = "1" ] && [ -z "$DASHBOARD_PASSWORD" ]; then
  DASHBOARD_PASSWORD="$(random_token | cut -c1-20)"
fi

{
  echo "bindAddr = \"$SERVER_ADDR\""
  echo "bindPort = $BIND_PORT"
  echo
  echo "[auth]"
  echo "method = \"token\""
  echo "token = \"$TOKEN\""

  if [ "$DASHBOARD" = "1" ]; then
    echo
    echo "[webServer]"
    echo "addr = \"$DASHBOARD_ADDR\""
    echo "port = $DASHBOARD_PORT"
    echo "user = \"$DASHBOARD_USER\""
    echo "password = \"$DASHBOARD_PASSWORD\""
  fi
} >"$OUTPUT"

echo "Wrote: $OUTPUT"
echo "token: $TOKEN"
if [ "$DASHBOARD" = "1" ]; then
  echo "dashboard: http://$DASHBOARD_ADDR:$DASHBOARD_PORT"
  echo "dashboard user: $DASHBOARD_USER"
  echo "dashboard password: $DASHBOARD_PASSWORD"
fi
