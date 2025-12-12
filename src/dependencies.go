package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// AnalyzeDependencies analyzes dependencies across different languages
func (c *Crawler) AnalyzeDependencies() {
	c.analyzePackageManagers()
	c.analyzeImports()
}

// analyzePackageManagers detects and parses package manager files
func (c *Crawler) analyzePackageManagers() {
	packageFiles := map[string]string{
		"package.json":     "npm",
		"requirements.txt": "pip",
		"Pipfile":          "pipenv",
		"poetry.lock":      "poetry",
		"go.mod":           "go modules",
		"Cargo.toml":       "cargo",
		"composer.json":    "composer",
		"Gemfile":          "bundler",
		"pom.xml":          "maven",
		"build.gradle":     "gradle",
		"Package.swift":    "swift pm",
		"pubspec.yaml":     "pub",
	}

	for _, files := range c.Analysis.FilesByType {
		for _, file := range files {
			baseName := filepath.Base(file.Path)
			if pmName, exists := packageFiles[baseName]; exists {
				c.parsePackageFile(file.Path, pmName)
			}
		}
	}
}

// parsePackageFile parses a package manager file
func (c *Crawler) parsePackageFile(filePath, pmName string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	pm := &PackageManager{
		Name:         pmName,
		ConfigFiles:  []string{filePath},
		Dependencies: make(map[string]string),
		DevDeps:      make(map[string]string),
	}

	switch pmName {
	case "npm":
		c.parsePackageJSON(content, pm)
	case "pip", "pipenv":
		c.parseRequirementsTxt(content, pm)
	case "go modules":
		c.parseGoMod(content, pm)
	case "cargo":
		c.parseCargoToml(content, pm)
	default:
		// Generic parsing
		c.parseGenericDeps(content, pm)
	}

	if len(pm.Dependencies) > 0 || len(pm.DevDeps) > 0 {
		c.Analysis.Dependencies.PackageManagers[pmName] = pm

		// Collect external dependencies
		for dep := range pm.Dependencies {
			c.Analysis.Dependencies.ExternalDeps = append(c.Analysis.Dependencies.ExternalDeps, dep)
		}
	}
}

// parsePackageJSON parses package.json
func (c *Crawler) parsePackageJSON(content []byte, pm *PackageManager) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return
	}

	if deps, ok := data["dependencies"].(map[string]interface{}); ok {
		for name, version := range deps {
			if v, ok := version.(string); ok {
				pm.Dependencies[name] = v
			}
		}
	}

	if devDeps, ok := data["devDependencies"].(map[string]interface{}); ok {
		for name, version := range devDeps {
			if v, ok := version.(string); ok {
				pm.DevDeps[name] = v
			}
		}
	}
}

// parseRequirementsTxt parses requirements.txt
func (c *Crawler) parseRequirementsTxt(content []byte, pm *PackageManager) {
	lines := strings.Split(string(content), "\n")
	re := regexp.MustCompile(`^([a-zA-Z0-9\-_]+)([>=<~!]+.*)?$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) >= 2 {
			name := matches[1]
			version := "unspecified"
			if len(matches) >= 3 && matches[2] != "" {
				version = matches[2]
			}
			pm.Dependencies[name] = version
		}
	}
}

// parseGoMod parses go.mod
func (c *Crawler) parseGoMod(content []byte, pm *PackageManager) {
	lines := strings.Split(string(content), "\n")
	inRequire := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		}

		if inRequire && line == ")" {
			inRequire = false
			continue
		}

		if inRequire || strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[0]
				if name == "require" && len(parts) >= 3 {
					name = parts[1]
					pm.Dependencies[name] = parts[2]
				} else {
					version := parts[1]
					pm.Dependencies[name] = version
				}
			}
		}
	}
}

// parseCargoToml parses Cargo.toml
func (c *Crawler) parseCargoToml(content []byte, pm *PackageManager) {
	lines := strings.Split(string(content), "\n")
	inDeps := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "[dependencies]" {
			inDeps = true
			continue
		}

		if strings.HasPrefix(line, "[") && line != "[dependencies]" {
			inDeps = false
		}

		if inDeps && strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				version := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				pm.Dependencies[name] = version
			}
		}
	}
}

// parseGenericDeps attempts generic dependency parsing
func (c *Crawler) parseGenericDeps(content []byte, pm *PackageManager) {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		pm.Dependencies[line] = "unknown"
	}
}

// analyzeImports analyzes import statements in source files
func (c *Crawler) analyzeImports() {
	importPatterns := map[string]*regexp.Regexp{
		"Python":     regexp.MustCompile(`(?m)^(?:from\s+([\w\.]+)|import\s+([\w\.]+))`),
		"JavaScript": regexp.MustCompile(`(?m)^(?:import.*from\s+['"]([^'"]+)['"]|require\(['"]([^'"]+)['"]\))`),
		"Go":         regexp.MustCompile(`(?m)^import\s+(?:[\w]+\s+)?["']([^"']+)["']`),
		"Rust":       regexp.MustCompile(`(?m)^use\s+([\w:]+)`),
		"Java":       regexp.MustCompile(`(?m)^import\s+([\w\.]+)`),
		"C#":         regexp.MustCompile(`(?m)^using\s+([\w\.]+)`),
		"Ruby":       regexp.MustCompile(`(?m)^require\s+['"]([^'"]+)['"]`),
	}

	for lang, files := range c.Analysis.FilesByType {
		pattern, exists := importPatterns[lang]
		if !exists {
			continue
		}

		for _, file := range files {
			content, err := os.ReadFile(file.Path)
			if err != nil {
				continue
			}

			matches := pattern.FindAllStringSubmatch(string(content), -1)
			var imports []string

			for _, match := range matches {
				for i := 1; i < len(match); i++ {
					if match[i] != "" {
						imports = append(imports, match[i])
					}
				}
			}

			if len(imports) > 0 {
				c.Analysis.Dependencies.ImportGraph[file.Path] = imports
			}
		}
	}
}
