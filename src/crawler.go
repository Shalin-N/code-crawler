package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CrawlerConfig holds configuration for the crawler
type CrawlerConfig struct {
	TargetPath  string
	OutputPath  string
	ExcludeDirs []string
	Verbose     bool
}

// Crawler is the main crawler instance
type Crawler struct {
	Config   *CrawlerConfig
	Analysis *Analysis
}

// Analysis holds all the collected data
type Analysis struct {
	RepoPath     string                `json:"repo_path"`
	AnalyzedAt   time.Time             `json:"analyzed_at"`
	Summary      *Summary              `json:"summary"`
	FileTree     *FileNode             `json:"file_tree"`
	FilesByType  map[string][]FileInfo `json:"files_by_type"`
	Dependencies *DependencyAnalysis   `json:"dependencies"`
	Statistics   *Statistics           `json:"statistics"`
}

// Summary provides high-level overview
type Summary struct {
	TotalFiles   int            `json:"total_files"`
	TotalDirs    int            `json:"total_dirs"`
	TotalSize    int64          `json:"total_size"`
	Languages    map[string]int `json:"languages"`
	LargestFiles []FileInfo     `json:"largest_files"`
	DeepestPath  string         `json:"deepest_path"`
	MaxDepth     int            `json:"max_depth"`
}

// FileNode represents a node in the file tree
type FileNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"is_dir"`
	Size     int64       `json:"size"`
	Children []*FileNode `json:"children,omitempty"`
}

// FileInfo holds information about a file
type FileInfo struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Extension string `json:"extension"`
	Language  string `json:"language"`
	Lines     int    `json:"lines,omitempty"`
}

// DependencyAnalysis holds dependency information
type DependencyAnalysis struct {
	PackageManagers map[string]*PackageManager `json:"package_managers"`
	ImportGraph     map[string][]string        `json:"import_graph"`
	ExternalDeps    []string                   `json:"external_deps"`
}

// PackageManager represents a detected package manager
type PackageManager struct {
	Name         string            `json:"name"`
	ConfigFiles  []string          `json:"config_files"`
	Dependencies map[string]string `json:"dependencies"`
	DevDeps      map[string]string `json:"dev_dependencies,omitempty"`
}

// Statistics provides detailed statistics
type Statistics struct {
	TotalLines      int            `json:"total_lines"`
	CodeLines       int            `json:"code_lines"`
	CommentLines    int            `json:"comment_lines"`
	BlankLines      int            `json:"blank_lines"`
	AvgFileSize     int64          `json:"avg_file_size"`
	FilesByLanguage map[string]int `json:"files_by_language"`
}

// NewCrawler creates a new crawler instance
func NewCrawler(config *CrawlerConfig) *Crawler {
	return &Crawler{
		Config: config,
		Analysis: &Analysis{
			RepoPath:    config.TargetPath,
			AnalyzedAt:  time.Now(),
			FilesByType: make(map[string][]FileInfo),
			Dependencies: &DependencyAnalysis{
				PackageManagers: make(map[string]*PackageManager),
				ImportGraph:     make(map[string][]string),
				ExternalDeps:    []string{},
			},
			Summary: &Summary{
				Languages: make(map[string]int),
			},
			Statistics: &Statistics{
				FilesByLanguage: make(map[string]int),
			},
		},
	}
}

// Scan performs the repository scan
func (c *Crawler) Scan() error {
	// Build file tree
	root, err := c.buildFileTree(c.Config.TargetPath, 0)
	if err != nil {
		return err
	}
	c.Analysis.FileTree = root

	// Collect file information
	if err := c.collectFileInfo(c.Config.TargetPath, 0); err != nil {
		return err
	}

	// Calculate summary statistics
	c.calculateSummary()

	return nil
}

// buildFileTree recursively builds the file tree structure
func (c *Crawler) buildFileTree(path string, depth int) (*FileNode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	node := &FileNode{
		Name:  filepath.Base(path),
		Path:  path,
		IsDir: info.IsDir(),
		Size:  info.Size(),
	}

	if !info.IsDir() {
		return node, nil
	}

	// Check if directory should be excluded
	if c.shouldExclude(filepath.Base(path)) {
		return node, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		// Skip directories we can't read
		return node, nil
	}

	for _, entry := range entries {
		if c.shouldExclude(entry.Name()) {
			continue
		}

		childPath := filepath.Join(path, entry.Name())
		childNode, err := c.buildFileTree(childPath, depth+1)
		if err != nil {
			continue // Skip problematic files
		}

		node.Children = append(node.Children, childNode)
	}

	return node, nil
}

// collectFileInfo collects detailed information about files
func (c *Crawler) collectFileInfo(path string, depth int) error {
	if depth > c.Analysis.Summary.MaxDepth {
		c.Analysis.Summary.MaxDepth = depth
		c.Analysis.Summary.DeepestPath = path
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil // Skip files we can't access
	}

	if info.IsDir() {
		c.Analysis.Summary.TotalDirs++

		if c.shouldExclude(filepath.Base(path)) {
			return nil
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return nil // Skip directories we can't read
		}

		for _, entry := range entries {
			if c.shouldExclude(entry.Name()) {
				continue
			}
			childPath := filepath.Join(path, entry.Name())
			c.collectFileInfo(childPath, depth+1)
		}
	} else {
		c.Analysis.Summary.TotalFiles++
		c.Analysis.Summary.TotalSize += info.Size()

		ext := filepath.Ext(path)
		lang := detectLanguage(ext, filepath.Base(path))

		fileInfo := FileInfo{
			Path:      path,
			Name:      info.Name(),
			Size:      info.Size(),
			Extension: ext,
			Language:  lang,
		}

		// Count lines if it's a text file
		if isTextFile(ext) {
			lines, err := countLines(path)
			if err == nil {
				fileInfo.Lines = lines
				c.Analysis.Statistics.TotalLines += lines
			}
		}

		c.Analysis.FilesByType[lang] = append(c.Analysis.FilesByType[lang], fileInfo)
		c.Analysis.Summary.Languages[lang]++
		c.Analysis.Statistics.FilesByLanguage[lang]++
	}

	return nil
}

// shouldExclude checks if a directory/file should be excluded
func (c *Crawler) shouldExclude(name string) bool {
	for _, excluded := range c.Config.ExcludeDirs {
		if name == excluded || strings.HasPrefix(name, ".") && excluded == ".hidden" {
			return true
		}
	}
	return false
}

// calculateSummary calculates summary statistics
func (c *Crawler) calculateSummary() {
	// Find largest files
	var allFiles []FileInfo
	for _, files := range c.Analysis.FilesByType {
		allFiles = append(allFiles, files...)
	}

	// Sort and get top 10 largest files
	if len(allFiles) > 0 {
		// Simple bubble sort for top 10
		for i := 0; i < len(allFiles) && i < 10; i++ {
			for j := i + 1; j < len(allFiles); j++ {
				if allFiles[j].Size > allFiles[i].Size {
					allFiles[i], allFiles[j] = allFiles[j], allFiles[i]
				}
			}
		}
		if len(allFiles) > 10 {
			c.Analysis.Summary.LargestFiles = allFiles[:10]
		} else {
			c.Analysis.Summary.LargestFiles = allFiles
		}
	}

	// Calculate average file size
	if c.Analysis.Summary.TotalFiles > 0 {
		c.Analysis.Statistics.AvgFileSize = c.Analysis.Summary.TotalSize / int64(c.Analysis.Summary.TotalFiles)
	}
}

// SaveData saves the analysis data to JSON
func (c *Crawler) SaveData() error {
	// Create output directory
	if err := os.MkdirAll(c.Config.OutputPath, 0755); err != nil {
		return err
	}

	// Save JSON data
	dataPath := filepath.Join(c.Config.OutputPath, "analysis.json")
	data, err := json.MarshalIndent(c.Analysis, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataPath, data, 0644)
}
