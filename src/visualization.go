package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

//go:embed visualizations/templates/*.html
var templatesFS embed.FS

// GenerateVisualization creates an interactive HTML visualization of the analysis
func (c *Crawler) GenerateVisualization() error {
	outputPath := filepath.Join(c.Config.OutputPath, "visualization.html")

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create visualization file: %w", err)
	}
	defer f.Close()

	// Create function map with helper functions
	funcMap := template.FuncMap{
		"formatBytes": formatBytes,
		"toJSON":      toJSON,
		"relPath":     relPath,
	}

	// Parse all template files
	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFS(
		templatesFS,
		"visualizations/templates/*.html",
	)
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	// Prepare data for template
	data := struct {
		RepoName string
		Analysis *Analysis
	}{
		RepoName: filepath.Base(c.Config.TargetPath),
		Analysis: c.Analysis,
	}

	// Execute the index template
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	fmt.Printf("Visualization saved to: %s\n", outputPath)
	return nil
}

// toJSON marshals any interface to JSON string
func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// relPath returns a relative path for display
func relPath(fullPath string) string {
	if fullPath == "" {
		return ""
	}
	return filepath.Base(fullPath)
}
