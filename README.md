# Pathfinder

Cross‑platform desktop file searcher built with Wails and Go.

---

## 📥 Download

Grab the right build for your OS—no cloning or building needed:

- **Windows (64‑bit)**  
  [Download pathfinder.exe](https://github.com/dylan0804/pathfinder/releases/latest/download/pathfinder.exe)

- **macOS (Intel & Apple Silicon)**  
  [Download pathfinder.app.zip](https://github.com/dylan0804/pathfinder/releases/latest/download/pathfinder.app.zip)

- **Linux (64‑bit)**  
  [Download pathfinder-linux-amd64](https://github.com/dylan0804/pathfinder/releases/latest/download/pathfinder-linux-amd64)

Unzip (if needed), make it executable (`chmod +x` on Linux/macOS), and run!

---

## 🚀 Features

- Fast, concurrent directory walk  
- Case‑insensitive substring & regex search  
- Skips dot‑dirs (`.git`, `node_modules`, `vendor`)  
- Emits results in real time via Wails events

## ⚙️ Prerequisites

- Go 1.18+  
- Node.js 16+ (for frontend assets)  
- Wails v2 CLI

## 🔧 Quickstart (Production)

```bash
# Build frontend
npm --prefix frontend install
npm --prefix frontend run build

# Build production binaries
export WAILS_ENV=production
wails build

# Run the binary for your platform:
#   macOS:  build/bin/pathfinder/pathfinder.app
#   Windows: build/bin/pathfinder/pathfinder.exe
#   Linux:   build/bin/pathfinder/pathfinder-linux-amd64
