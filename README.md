# Pathfinder: Fast File Search Application

A cross-platform desktop file search application built with [Wails](https://wails.app/) (Go + web frontend). It uses Go’s built-in filesystem APIs and a simple concurrent walker to deliver fast, real-time search results via a lightweight UI.

## Features

- **Concurrent traversal** using Go goroutines and channels
- **Pattern matching** (substring or regex)
- **Minimal dependencies**: single binary with embedded frontend assets
- **Cross-platform**: runs on Windows, macOS, and Linux

## Prerequisites

- Go **1.18+**
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

2. Build frontend assets:

   ```bash
   cd frontend
   npm install      # or yarn install
   npm run build    # outputs to frontend/dist/
   cd ..
   ```

3. Run in development mode (with live reload):

   ```bash
   wails dev
   ```

## Production Build

1. Export the production environment variable:

   ```bash
   # macOS/Linux (bash/zsh)
   export WAILS_ENV=production

   # Windows (CMD)
   set WAILS_ENV=production

   # Windows (PowerShell)
   $env:WAILS_ENV = "production"
   ```

2. Build the application:

   ```bash
   wails build
   ```

3. Run the binary:

   ```bash
   cd build/bin/pathfinder
   
   # macOS: open the .app bundle
   open pathfinder.app
   
   # Windows: run the .exe
   .\pathfinder.exe
   ```

## Cross-Compilation (Example)

From macOS/Linux terminal:

```bash
export WAILS_ENV=production
npm --prefix frontend install
npm --prefix frontend run build

# Cross-compile for Windows 64-bit:
GOOS=windows GOARCH=amd64 wails build
```

After building, the Windows executable is located in:

```
build/bin/pathfinder/pathfinder.exe
```

## License

MIT © Dylan Christian Dihalim

