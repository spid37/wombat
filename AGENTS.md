# AGENTS.md

Guidance for AI agents working on **Wombat** (cross-platform gRPC client).

## Stack

- **Backend:** Go 1.27+, module `wombat`
- **UI shell:** Wails v2 (WKWebView / WebKitGTK / WebView2)
- **Frontend:** Svelte 3 + Rollup 4 in `frontend/`
- **Bindings:** generated under `frontend/wailsjs/` (`wails generate module` / `wails build`)

## Architecture

- Bound Go API: `internal/app.API` — public methods are exposed to JS
- Frontend calls Go via `window.backend.api.*` (`frontend/src/wails-bridge.js`)
- Events use `frontend/src/runtime.js` (`Events.On` / `Events.Emit`) — **never overwrite `window.wails`** (Wails owns it for `EventsNotify`)
- Event name constants: `internal/app/events.go`
- Flux-style: frontend mostly reacts to events, not returned promises
- Dev/debug builds start a test gRPC server on `localhost:5001` (`internal/server`)

## Commands

```zsh
go install github.com/wailsapp/wails/v2/cmd/wails@latest

make help
make dev                 # live reload
make build               # native host build
make build-darwin        # macOS arm64 (sets macOS 13 CGO flags)
make build-linux
make build-windows
```

Release CI runs on push of tags matching `v*` (see `.github/workflows/package.yml`).

## Critical pitfalls

1. **Assets:** `main.go` embeds `frontend/public` then `fs.Sub(..., "frontend/public")`. Do not pass the raw embed FS to Wails or the UI will be blank.
2. **`window.wails`:** Do not replace it with a compatibility object. Use `frontend/src/runtime.js` for Events/Browser.
3. **macOS linker warning:** Go 1.27 builds for macOS 13; bare `wails build` may inject CGO min 10.13. Prefer `make build-darwin` (or set `MACOSX_DEPLOYMENT_TARGET` / `CGO_CFLAGS` / `CGO_LDFLAGS` to `13.0`).
4. **Monaco:** Import from `frontend/src/monaco.js`. Stay on Monaco ~0.52 unless import paths are updated (0.56 broke ESM paths).

## Conventions

- Keep changes focused; match existing style in `internal/app` and `frontend/src/views`
- Prefer events from `events.go` over inventing new frontend/backend protocols
- After frontend API/binding changes, regenerate with `wails generate module` or a full `wails build`
- Note user-facing changes in `CHANGELOG.md` under `[Unreleased]`
- Do not commit `build/bin/` or secrets; keep `build/appicon.png` and `build/darwin/Info.plist`

## Out of scope unless requested

- Badger v2 → v4 migration
- Svelte 5 / major UI rewrites
- Apple notarization / code signing
