# Visualizations Module

This directory contains the HTML visualization components for the code crawler.

## Structure

```
visualizations/
└── templates/           # HTML template components
    ├── layout.html      # Main layout template
    ├── styles.html      # CSS styles
    ├── header.html      # Page header
    ├── stats.html       # Statistics cards
    ├── charts.html      # Chart containers
    ├── languages.html   # Language badges
    ├── files.html       # Largest files list
    ├── dependencies.html # Dependencies section
    ├── tree.html        # Directory tree
    ├── scripts.html     # JavaScript code
    └── footer.html      # Page footer
```

## How It Works

The visualization system uses Go's `html/template` package with embedded templates:

1. **Template Embedding**: All HTML templates are embedded into the binary using `//go:embed`
2. **Component-Based**: Each section is a separate template file for easy maintenance
3. **Main Layout**: `layout.html` imports and orchestrates all components
4. **No External Files**: Templates are compiled into the binary, making it portable

## Adding New Components

1. Create a new template file in `templates/`:
```html
{{define "my-component"}}
<div class="my-component">
    <!-- Your HTML here -->
</div>
{{end}}
```

2. Include it in `layout.html`:
```html
{{template "my-component" .}}
```

## Modifying Existing Components

Simply edit the relevant template file:

- **Styles**: Edit `styles.html`
- **Charts**: Edit `charts.html` and `scripts.html`
- **Layout**: Edit `layout.html`

Changes will be reflected after rebuilding the binary.

## Template Data

All templates receive the same data structure:

```go
{
    RepoName: string           // Repository name
    Analysis: *Analysis {      // Analysis results
        Summary: {
            TotalFiles, TotalDirs, TotalSize, MaxDepth
            Languages: map[string]int
            LargestFiles: []FileInfo
        }
        Statistics: {
            TotalLines, CodeLines, etc.
        }
        Dependencies: {
            PackageManagers: map[string]*PackageManager
        }
        FileTree: *FileNode
    }
}
```

## Helper Functions

Available in all templates:

- `formatBytes`: Format byte sizes (e.g., "1.5 MB")
- `toJSON`: Convert data to JSON string
- `relPath`: Get relative path for display

## Development

After modifying templates, rebuild the project:

```bash
go build -o code-crawler ./src
```

The templates are embedded at build time, so changes require a rebuild.
