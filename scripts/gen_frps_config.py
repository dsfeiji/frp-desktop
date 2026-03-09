#!/usr/bin/env python3
"""Generate frps.toml for Linux/Windows servers.

Usage (Linux/macOS):
  python3 scripts/gen_frps_config.py --server-addr 0.0.0.0 --bind-port 7000 --token your_token

Usage (Windows):
  py scripts\\gen_frps_config.py --server-addr 0.0.0.0 --bind-port 7000 --token your_token
"""

from __future__ import annotations

import argparse
import pathlib
import secrets
import string


def random_token(length: int = 32) -> str:
    alphabet = string.ascii_letters + string.digits
    return "".join(secrets.choice(alphabet) for _ in range(length))


def build_toml(args: argparse.Namespace) -> str:
    lines: list[str] = []
    lines.append(f'bindAddr = "{args.server_addr}"')
    lines.append(f"bindPort = {args.bind_port}")
    lines.append("")
    lines.append("[auth]")
    lines.append('method = "token"')
    lines.append(f'token = "{args.token}"')

    if args.dashboard:
        lines.append("")
        lines.append("[webServer]")
        lines.append(f'addr = "{args.dashboard_addr}"')
        lines.append(f"port = {args.dashboard_port}")
        lines.append(f'user = "{args.dashboard_user}"')
        lines.append(f'password = "{args.dashboard_password}"')

    return "\n".join(lines) + "\n"


def main() -> None:
    parser = argparse.ArgumentParser(description="Generate frps.toml for FRP server")
    parser.add_argument("--server-addr", default="0.0.0.0", help="frps bind address")
    parser.add_argument("--bind-port", type=int, default=7000, help="frps bind port")
    parser.add_argument("--token", default="", help="auth token shared with clients")
    parser.add_argument("--output", default="frps.toml", help="output config path")

    parser.add_argument("--dashboard", action="store_true", help="enable dashboard section")
    parser.add_argument("--dashboard-addr", default="0.0.0.0", help="dashboard bind address")
    parser.add_argument("--dashboard-port", type=int, default=7500, help="dashboard port")
    parser.add_argument("--dashboard-user", default="admin", help="dashboard username")
    parser.add_argument("--dashboard-password", default="", help="dashboard password")

    args = parser.parse_args()

    if not (1 <= args.bind_port <= 65535):
        raise SystemExit("bind-port must be 1-65535")

    if args.dashboard and not (1 <= args.dashboard_port <= 65535):
        raise SystemExit("dashboard-port must be 1-65535")

    if not args.token:
        args.token = random_token(32)

    if args.dashboard and not args.dashboard_password:
        args.dashboard_password = random_token(20)

    output = pathlib.Path(args.output)
    output.write_text(build_toml(args), encoding="utf-8")

    print(f"Wrote: {output.resolve()}")
    print(f"token: {args.token}")
    if args.dashboard:
        print(f"dashboard: http://{args.dashboard_addr}:{args.dashboard_port}")
        print(f"dashboard user: {args.dashboard_user}")
        print(f"dashboard password: {args.dashboard_password}")


if __name__ == "__main__":
    main()
