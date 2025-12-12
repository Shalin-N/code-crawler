package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const version = "1.0.0"

func main() {
	// CLI flags
	targetPath := flag.String("path", ".", "Path to the repository to analyze")
	outputPath := flag.String("output", ".analysis", "Output directory for analysis files")
	excludeDirs := flag.String("exclude", ".git,node_modules,vendor,.dist,build,target,.venv,__pycache__", "Comma-separated list of directories to exclude")
	generateViz := flag.Bool("viz", true, "Generate HTML visualization")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	showVersion := flag.Bool("version", false, "Show version information")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Code Crawler v%s\n", version)
		os.Exit(0)
	}

	// Validate target path
	absPath, err := filepath.Abs(*targetPath)
	if err != nil {
		log.Fatalf("Invalid path: %v", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		log.Fatalf("Path does not exist: %s", absPath)
	}

	fmt.Printf("🔍 Code Crawler v%s\n", version)
	fmt.Printf("📂 Analyzing: %s\n", absPath)
	fmt.Println(strings.Repeat("=", 60))

	startTime := time.Now()

	// Initialize crawler
	config := &CrawlerConfig{
		TargetPath:  absPath,
		OutputPath:  *outputPath,
		ExcludeDirs: parseExcludeDirs(*excludeDirs),
		Verbose:     *verbose,
	}

	crawler := NewCrawler(config)

	// Scan the repository
	fmt.Println("\n📊 Scanning repository structure...")
	if err := crawler.Scan(); err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	// Analyze dependencies
	fmt.Println("\n🔗 Analyzing dependencies...")
	crawler.AnalyzeDependencies()

	// Save analysis data
	fmt.Println("\n💾 Saving analysis data...")
	if err := crawler.SaveData(); err != nil {
		log.Fatalf("Failed to save data: %v", err)
	}

	// Generate visualizations
	if *generateViz {
		fmt.Println("\n🎨 Generating visualizations...")
		if err := crawler.GenerateVisualization(); err != nil {
			log.Fatalf("Failed to generate visualization: %v", err)
		}
	}

	duration := time.Since(startTime)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\n✅ Analysis complete in %.2f seconds\n", duration.Seconds())
	fmt.Printf("📁 Output saved to: %s\n", *outputPath)
	if *generateViz {
		vizPath := filepath.Join(*outputPath, "visualization.html")
		fmt.Printf("🌐 Open visualization: file://%s\n", vizPath)
	}
}

func parseExcludeDirs(excludeStr string) []string {
	if excludeStr == "" {
		return []string{}
	}
	dirs := strings.Split(excludeStr, ",")
	for i := range dirs {
		dirs[i] = strings.TrimSpace(dirs[i])
	}
	return dirs
}
