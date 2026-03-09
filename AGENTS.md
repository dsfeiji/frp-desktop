# Repository Guidelines

## Project Structure & Module Organization
- Root Go app entrypoints are `main.go` and `app.go` (Wails backend + process/control logic).
- Frontend lives in `frontend/`:
  - `src/main.js` for UI behavior and Wails bindings usage
  - `src/style.css` / `src/app.css` for styling
  - `dist/` is generated build output (do not hand-edit)
- Packaging/build assets are under `build/` (platform manifests, icons, installer files).
- Release artifacts are under `release/` and should be treated as generated output.
- Ignore duplicate scratch files like `app (1).go` unless intentionally used.

## Build, Test, and Development Commands
- `wails dev`: run desktop app with live frontend/backend development.
- `wails build`: produce production desktop build using Wails config.
- `go build ./...`: compile all Go packages; quick backend sanity check.
- `npm --prefix frontend run dev`: run Vite frontend dev server only.
- `npm --prefix frontend run build`: build frontend static assets.
- `npm --prefix frontend run preview`: preview built frontend output.

## Coding Style & Naming Conventions
- Go: format with `gofmt` (tabs, idiomatic Go naming, exported `CamelCase`, internal `camelCase`).
- JavaScript/CSS: 2-space indentation; keep functions small and event handlers explicit.
- Prefer descriptive names tied to FRP domain (`StartFrpc`, `LocalPorts`, `RuntimeState`).
- Keep constants centralized (see fixed server/auth values in `app.go`) and avoid magic numbers.

## Testing Guidelines
- No automated tests are currently present (`*_test.go` not found).
- Minimum pre-PR verification:
  - `go build ./...`
  - `npm --prefix frontend run build`
  - `wails build` (or `wails build -s` for local release checks)
- When adding tests, use Go’s `testing` package with file pattern `*_test.go` and table-driven cases for port/config parsing logic.

## Commit & Pull Request Guidelines
- `.git` metadata is not available in this snapshot, so existing commit history conventions cannot be inferred here.
- Use Conventional Commit style going forward (e.g., `feat: add port validation`, `fix: handle frpc path errors`).
- PRs should include:
  - Clear scope and user-facing impact
  - Verification steps and command output summary
  - Screenshots/GIFs for UI changes in `frontend/src`
  - Linked issue/task when applicable

## Security & Configuration Notes
- Do not commit real secrets/tokens. Current fixed auth/server values in source should be moved to configurable secure inputs before production use.
