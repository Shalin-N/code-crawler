# 🚀 Quick Start Guide

## Installation

### Build from Source

```bash
# Clone the repository
git clone https://github.com/shalinnaidoo/code-crawler.git
cd code-crawler

# Build the executable
go build -o code-crawler

# Or use the build script for all platforms
chmod +x build.sh
./build.sh
```

### Download Binary

Download the pre-built binary for your platform from the [Releases](https://github.com/shalinnaidoo/code-crawler/releases) page.

## Basic Usage

```bash
# Analyze current directory
./code-crawler

# Analyze specific path
./code-crawler -path /path/to/your/repo

# Custom output location
./code-crawler -path ~/projects/my-app -output ~/Desktop/analysis

# Skip visualization generation
./code-crawler -viz=false

# Exclude specific directories
./code-crawler -exclude ".git,node_modules,vendor,build"
```

## Example Output

After running the crawler, you'll get:

1. **analysis.json** - Complete analysis data in JSON format
2. **visualization.html** - Interactive HTML dashboard

Simply open `visualization.html` in your browser to explore your codebase!

## Testing the Tool

Test it on this repository itself:

```bash
./code-crawler -path . -output ./self-analysis
```

Then open `self-analysis/visualization.html` in your browser.

## Building for Multiple Platforms

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o code-crawler-mac-arm64

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o code-crawler-mac-amd64

# Linux
GOOS=linux GOARCH=amd64 go build -o code-crawler-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o code-crawler.exe

# Or use the build script
./build.sh
```

## What Gets Analyzed?

- 📁 **Directory Structure** - Complete file tree
- 🔤 **Languages** - 40+ languages detected automatically
- 📦 **Dependencies** - Package managers (npm, pip, go mod, cargo, maven, etc.)
- 📊 **Statistics** - File counts, sizes, line counts
- 🔗 **Import Graph** - How files import each other

## Tips

- Run from the root of your repository
- The tool automatically excludes common build directories (.git, node_modules, etc.)
- For large repos, analysis completes in seconds
- All data is stored locally - nothing leaves your machine
- The visualization works offline after generation

## Example Projects to Try

```bash
# Analyze a React project
./code-crawler -path ~/projects/my-react-app

# Analyze a Python project
./code-crawler -path ~/projects/django-blog

# Analyze a Go project
./code-crawler -path ~/projects/go-service

# Analyze a Rust project
./code-crawler -path ~/projects/rust-cli
```

Enjoy exploring your codebase! 🎉
