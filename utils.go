package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Language detection mappings
var languageMap = map[string]string{
	// Programming Languages
	".go":    "Go",
	".py":    "Python",
	".js":    "JavaScript",
	".ts":    "TypeScript",
	".jsx":   "JavaScript",
	".tsx":   "TypeScript",
	".java":  "Java",
	".c":     "C",
	".cpp":   "C++",
	".cc":    "C++",
	".cxx":   "C++",
	".h":     "C/C++ Header",
	".hpp":   "C++ Header",
	".cs":    "C#",
	".rb":    "Ruby",
	".php":   "PHP",
	".swift": "Swift",
	".kt":    "Kotlin",
	".rs":    "Rust",
	".scala": "Scala",
	".r":     "R",
	".m":     "Objective-C",
	".dart":  "Dart",
	".lua":   "Lua",
	".pl":    "Perl",
	".sh":    "Shell",
	".bash":  "Shell",
	".zsh":   "Shell",
	".fish":  "Shell",
	".sql":   "SQL",

	// Web
	".html":   "HTML",
	".htm":    "HTML",
	".css":    "CSS",
	".scss":   "SCSS",
	".sass":   "Sass",
	".less":   "Less",
	".vue":    "Vue",
	".svelte": "Svelte",

	// Data & Config
	".json": "JSON",
	".yaml": "YAML",
	".yml":  "YAML",
	".xml":  "XML",
	".toml": "TOML",
	".ini":  "INI",
	".env":  "Environment",
	".conf": "Config",
	".cfg":  "Config",

	// Documentation
	".md":   "Markdown",
	".rst":  "reStructuredText",
	".txt":  "Text",
	".tex":  "LaTeX",
	".adoc": "AsciiDoc",

	// Build & Package
	".gradle":     "Gradle",
	".maven":      "Maven",
	".dockerfile": "Dockerfile",
	".mk":         "Makefile",

	// Others
	".proto":   "Protocol Buffers",
	".graphql": "GraphQL",
	".gql":     "GraphQL",
}

// Special filenames
var specialFiles = map[string]string{
	"Dockerfile":        "Dockerfile",
	"Makefile":          "Makefile",
	"Rakefile":          "Ruby",
	"Gemfile":           "Ruby",
	"Podfile":           "Ruby",
	"CMakeLists.txt":    "CMake",
	"package.json":      "JSON",
	"tsconfig.json":     "JSON",
	"webpack.config.js": "JavaScript",
	"rollup.config.js":  "JavaScript",
	"vite.config.js":    "JavaScript",
	"vue.config.js":     "JavaScript",
	".gitignore":        "Config",
	".dockerignore":     "Config",
	".eslintrc":         "JSON",
	".prettierrc":       "JSON",
	"requirements.txt":  "Text",
	"go.mod":            "Go Module",
	"go.sum":            "Go Module",
	"Cargo.toml":        "TOML",
	"Cargo.lock":        "TOML",
	"pyproject.toml":    "TOML",
	"Pipfile":           "TOML",
	"pom.xml":           "XML",
	"build.gradle":      "Gradle",
}

// detectLanguage detects the programming language from file extension
func detectLanguage(ext, filename string) string {
	// Check special filenames first
	if lang, exists := specialFiles[filename]; exists {
		return lang
	}

	// Check by extension
	ext = strings.ToLower(ext)
	if lang, exists := languageMap[ext]; exists {
		return lang
	}

	// Check if it's hidden file
	if strings.HasPrefix(filename, ".") && ext == "" {
		return "Config"
	}

	// Unknown
	if ext == "" {
		return "No Extension"
	}
	return "Other"
}

// isTextFile checks if a file is likely a text file
func isTextFile(ext string) bool {
	textExtensions := map[string]bool{
		".go": true, ".py": true, ".js": true, ".ts": true,
		".jsx": true, ".tsx": true, ".java": true, ".c": true,
		".cpp": true, ".cc": true, ".cxx": true, ".h": true,
		".hpp": true, ".cs": true, ".rb": true, ".php": true,
		".swift": true, ".kt": true, ".rs": true, ".scala": true,
		".r": true, ".m": true, ".dart": true, ".lua": true,
		".pl": true, ".sh": true, ".bash": true, ".zsh": true,
		".fish": true, ".sql": true, ".html": true, ".htm": true,
		".css": true, ".scss": true, ".sass": true, ".less": true,
		".vue": true, ".svelte": true, ".json": true, ".yaml": true,
		".yml": true, ".xml": true, ".toml": true, ".ini": true,
		".md": true, ".rst": true, ".txt": true, ".tex": true,
		".proto": true, ".graphql": true, ".gql": true,
	}

	return textExtensions[strings.ToLower(ext)]
}

// countLines counts the number of lines in a file
func countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

// formatBytes formats byte size into human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// relativePath returns a path relative to the base directory
func relativePath(basePath, targetPath string) string {
	rel, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return targetPath
	}
	return rel
}
