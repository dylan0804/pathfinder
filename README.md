# Pathfinder: Fast File Search Application

A cross-platform desktop file search application built with [Wails](https://wails.app/) (Go + web frontend). It uses Go’s built-in filesystem APIs and a simple concurrent walker to deliver fast, real-time search results via a lightweight UI.

## Features

- **Concurrent traversal** using Go goroutines and channels
- **Pattern matching** (substring or regex)
- **Minimal dependencies**: single binary with embedded frontend assets
- **Cross-compile** for Windows, macOS, and Linux

## Prerequisites

- Go **1.18+** installed
- Node.js **16+** and npm (or Yarn)
- [Wails CLI](https://wails.app/) v2:
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

## Getting Started

1. Clone the repo:
   ```bash
   git clone https://github.com/dylan0804/pathfinder.git
   cd pathfinder
   ```

2. Install frontend dependencies and build assets:
   ```bash
   cd frontend
   npm install      # or yarn
   npm run build    # outputs to frontend/dist/
   cd ..
   ```

3. Run in development mode (with live reload):
   ```bash
   wails dev
   ```

## Production Build
1. Set production environment:
   ```bash
   export WAILS_ENV=production      # macOS/Linux (bash/zsh)
   set WAILS_ENV=production         # Windows cmd
   $env:WAILS_ENV = "production"    # Windows PowerShell
   ```

2. Build the app:
   ```bash
   wails build
   ```

3. Run the binary:
   ```bash
   cd build/bin/pathfinder
   open pathfinder.app
   ```

## Cross Compilation

From macOS/Linux:
```bash
export WAILS_ENV=production
GOOS=windows GOARCH=amd64 wails build
```

The Windows executable will be in `build/bin/pathfinder.exe`.

## License

MIT © Dylan Christiandi Halim

