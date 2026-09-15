# Development

## Prerequisites 

- Working [Go](https://golang.org/) 1.27+ environment
- Working [Node.js](https://nodejs.org/en/) v20+ environment
- [Wails v2](https://wails.io/) CLI

## Getting started

Install the Wails CLI:

```zsh
$ go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Run in development mode (starts the Go backend and frontend watcher):

```zsh
$ wails dev
```

## Built-in gRPC server

During development it's common to want to test proto files and gRPC APIs. When Wombat launches in development/debug mode it
automatically starts a gRPC server running on `localhost:5001`. You can also make changes to the
`/internal/server/foobar.proto` file to add test protos.

## Backend API

The Go backend has a singular API entry: `app.API`; All public methods defined on the `API` struct are exposed to the
frontend via generated bindings in `frontend/wailsjs/go/app`.

Wombat follows a (loose) [flux](https://facebook.github.io/flux/) like architecture. API calls are made to the backend
(actions) and the backend will then emit events (dispatcher) for any changes required by the frontend.

Although the backend APIs can respond with an object/error; the frontend mostly ignores these objects/errors (uncaught promise error), and
rather responds to events including error events.

All backend event constants are defined in `/internal/app/events.go`.

## Frontend

The frontend is built with [Svelte](https://svelte.dev/); and should be fairly straight forward if you have done any
type of web-based frontend development before.

Wombat uses the system webview (macOS WKWebView, Linux WebKitGTK, Windows WebView2). A small bridge in
`frontend/src/wails-bridge.js` exposes `window.backend` / `window.wails` for the existing Svelte views.

If you need to add any styles please update both `frontend/src/views/App.svelte` and `frontend/rollup.config.js`.

After getting your UI correct; it's always good to validate in a real build:

```zsh
$ wails build -debug
$ open ./build/bin/Wombat.app   # macOS
```

In debug mode you can still "Right/Option Click" -> "Inspect Element" and bring up the development tools (mac and linux
only). This is disabled in production builds.

## Changelog

We maintain a changelog for all versions of Wombat. If you create a PR please add the change to `CHANGELOG.md`. Don't
forget to add your GitHub handle to make sure you get the credit!
