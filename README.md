# Code Crawler 🔍

A language-agnostic code analysis tool that helps you understand any codebase's architecture, dependencies, and structure. Compiles to a single executable that you can drop into any repository.

## Features

- 📂 **Directory Structure Analysis** - Complete file tree visualization
- 🔗 **Dependency Detection** - Supports npm, pip, Go modules, Cargo, Maven, and more
- 💻 **Language Detection** - Automatically identifies 40+ programming languages
- 📊 **Rich Visualizations** - Interactive HTML dashboard with charts
- 💾 **JSON Export** - All analysis data saved in structured format
- ⚡ **Fast & Lightweight** - Single executable, no dependencies required
- 🌍 **Language Agnostic** - Works with any programming language

## Installation

### Option 1: Download Pre-built Binary

*(Coming soon - check releases page)*

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/shalinnaidoo/code-crawler.git
cd code-crawler

# Build the executable
go build -o code-crawler

# Optional: Install to your PATH
go install
```

## Usage

### Basic Usage

```bash
# Analyze current directory
./code-crawler

# Analyze specific path
./code-crawler -path /path/to/repo

# Specify output directory
./code-crawler -path /path/to/repo -output ./analysis-output
```

### CLI Options

```bash
  -path string
        Path to the repository to analyze (default ".")
  
  -output string
        Output directory for analysis files (default "code-analysis")
  
  -exclude string
        Comma-separated list of directories to exclude 
        (default ".git,node_modules,vendor,dist,build,target,.venv,__pycache__")
  
  -viz
        Generate HTML visualization (default true)
  
  -verbose
        Enable verbose logging (default false)
  
  -version
        Show version information
```

### Examples

```bash
# Analyze a JavaScript project
./code-crawler -path ~/projects/my-react-app

# Analyze Python project with custom exclusions
./code-crawler -path ~/projects/django-app -exclude ".git,.venv,migrations"

# Quick analysis without visualization
./code-crawler -path ~/projects/go-api -viz=false

# Analyze and save to specific location
./code-crawler -path ~/projects/rust-cli -output ~/Desktop/analysis
```

## Output

Code Crawler generates two main outputs:

### 1. analysis.json

A comprehensive JSON file containing:
- File tree structure
- Files grouped by language/type
- Dependency information
- Statistics (line counts, file sizes, etc.)
- Language distribution

Example structure:
```json
{
  "repo_path": "/path/to/repo",
  "analyzed_at": "2025-12-13T10:30:00Z",
  "summary": {
    "total_files": 156,
    "total_dirs": 42,
    "total_size": 2458624,
    "languages": {
      "JavaScript": 45,
      "TypeScript": 38,
      "CSS": 12
    }
  },
  "dependencies": {
    "package_managers": {
      "npm": {
        "dependencies": {
          "react": "^18.2.0",
          "express": "^4.18.2"
        }
      }
    }
  }
}
```

### 2. visualization.html

An interactive HTML dashboard featuring:
- 📊 Visual statistics cards
- 🥧 Language distribution pie chart
- 📈 Files per language bar chart
- 📦 Largest files listing
- 📚 Dependencies breakdown
- 🌳 Interactive directory tree

Simply open the HTML file in your browser!

## Supported Languages

Code Crawler detects 40+ languages including:

**Programming Languages:**
- Go, Python, JavaScript, TypeScript, Java, C, C++, C#, Rust, Ruby, PHP, Swift, Kotlin, Scala, Dart, Lua, Perl, R, Objective-C

**Web Technologies:**
- HTML, CSS, SCSS, Sass, Less, Vue, Svelte

**Data & Config:**
- JSON, YAML, XML, TOML, INI, Environment files

**Documentation:**
- Markdown, reStructuredText, LaTeX, AsciiDoc

**Build & Package:**
- Dockerfile, Makefile, Gradle, Maven, CMake

**And many more...**

## Dependency Detection

Automatically detects and parses:

| Language/Framework | Package Manager | Config Files |
|-------------------|-----------------|--------------|
| JavaScript/Node.js | npm | package.json |
| Python | pip, pipenv, poetry | requirements.txt, Pipfile, poetry.lock |
| Go | Go modules | go.mod |
| Rust | Cargo | Cargo.toml |
| Java | Maven, Gradle | pom.xml, build.gradle |
| PHP | Composer | composer.json |
| Ruby | Bundler | Gemfile |
| Swift | Swift PM | Package.swift |
| Dart | Pub | pubspec.yaml |

## Use Cases

### 1. Onboarding to New Codebases
Quickly understand the structure and dependencies of an unfamiliar project.

### 2. Code Audits
Generate comprehensive reports for code reviews or audits.

### 3. Documentation
Create up-to-date visualizations of your project structure.

### 4. Migration Planning
Analyze dependencies before migrating to new frameworks or languages.

### 5. Technical Debt Analysis
Identify large files, deep nesting, and complexity hotspots.

## Building for Multiple Platforms

```bash
# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o code-crawler-mac-amd64

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o code-crawler-mac-arm64

# Linux
GOOS=linux GOARCH=amd64 go build -o code-crawler-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o code-crawler.exe
```

## Advanced Configuration

### Custom Exclusions

Create a `.code-crawler.json` config file in your project:

```json
{
  "exclude": [
    ".git",
    "node_modules",
    "coverage",
    "dist",
    "build"
  ],
  "max_file_size": 10485760,
  "include_hidden": false
}
```

## Performance

- Analyzes 10,000+ files in seconds
- Low memory footprint
- Handles large repositories (100k+ files)
- Skips binary files automatically

## Roadmap

- [ ] Git history analysis
- [ ] Code complexity metrics (cyclomatic complexity)
- [ ] Dead code detection
- [ ] Integration with CI/CD pipelines
- [ ] Plugin system for custom analyzers
- [ ] VSCode extension
- [ ] Comparison between different versions

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - feel free to use this tool in your projects!

## Author

Created by Shalin Naidoo

## Acknowledgments

- Built with Go for maximum portability
- Chart.js for beautiful visualizations
- Inspired by tools like tokei, cloc, and tree

---

**Star this repo if you find it useful!** ⭐
