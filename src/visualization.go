package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// GenerateVisualization creates an interactive HTML visualization of the analysis
func (c *Crawler) GenerateVisualization() error {
	outputPath := filepath.Join(c.Config.OutputPath, "visualization.html")

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create visualization file: %w", err)
	}
	defer f.Close()

	funcMap := template.FuncMap{
		"formatBytes": formatBytes,
		"toJSON":      toJSON,
		"relPath":     relPath,
	}

	tmpl, err := template.New("visualization").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		RepoName string
		Analysis *Analysis
	}{
		RepoName: filepath.Base(c.Config.TargetPath),
		Analysis: c.Analysis,
	}

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

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.RepoName}} - Code Analysis</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js"></script>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 2rem;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            overflow: hidden;
        }

        header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 3rem 2rem;
            text-align: center;
        }

        header h1 {
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
            text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        }

        header p {
            opacity: 0.9;
            font-size: 1.1rem;
        }

        .content {
            padding: 2rem;
        }

        .section {
            margin-bottom: 3rem;
        }

        .section-title {
            font-size: 1.8rem;
            color: #333;
            margin-bottom: 1.5rem;
            padding-bottom: 0.5rem;
            border-bottom: 3px solid #667eea;
        }

        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 1.5rem;
            margin-bottom: 3rem;
        }

        .stat-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 1.5rem;
            border-radius: 12px;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
            transition: transform 0.3s ease;
        }

        .stat-card:hover {
            transform: translateY(-5px);
        }

        .stat-label {
            font-size: 0.9rem;
            opacity: 0.9;
            margin-bottom: 0.5rem;
        }

        .stat-value {
            font-size: 2rem;
            font-weight: bold;
        }

        .charts-container {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
            gap: 2rem;
            margin-bottom: 3rem;
        }

        .chart-box {
            background: #f8f9fa;
            padding: 1.5rem;
            border-radius: 12px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }

        .chart-title {
            font-size: 1.2rem;
            color: #333;
            margin-bottom: 1rem;
            text-align: center;
        }

        .languages-badges {
            display: flex;
            flex-wrap: wrap;
            gap: 0.75rem;
            margin-bottom: 2rem;
        }

        .language-badge {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 0.5rem 1rem;
            border-radius: 20px;
            font-size: 0.9rem;
            box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
        }

        .file-list {
            background: #f8f9fa;
            border-radius: 12px;
            padding: 1.5rem;
            margin-bottom: 2rem;
        }

        .file-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 1rem;
            background: white;
            border-radius: 8px;
            margin-bottom: 0.75rem;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.05);
        }

        .file-name {
            font-family: 'Courier New', monospace;
            color: #667eea;
            font-weight: bold;
        }

        .file-size {
            color: #666;
            font-size: 0.9rem;
        }

        .dependencies-section {
            background: #f8f9fa;
            border-radius: 12px;
            padding: 1.5rem;
            margin-bottom: 2rem;
        }

        .dep-group {
            margin-bottom: 1.5rem;
        }

        .dep-title {
            font-size: 1.2rem;
            color: #667eea;
            margin-bottom: 0.75rem;
            font-weight: bold;
        }

        .dep-list {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
            gap: 0.5rem;
        }

        .dep-item {
            background: white;
            padding: 0.75rem;
            border-radius: 6px;
            font-family: 'Courier New', monospace;
            font-size: 0.85rem;
            box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        }

        .tree {
            background: #f8f9fa;
            border-radius: 12px;
            padding: 1.5rem;
            font-family: 'Courier New', monospace;
            font-size: 0.9rem;
            overflow-x: auto;
        }

        .tree-item {
            padding: 0.25rem 0;
            color: #333;
        }

        .tree-dir {
            color: #667eea;
            font-weight: bold;
        }

        .tree-file {
            color: #666;
        }

        footer {
            background: #f8f9fa;
            padding: 2rem;
            text-align: center;
            color: #666;
            border-top: 1px solid #e0e0e0;
        }

        @media (max-width: 768px) {
            body {
                padding: 1rem;
            }

            header h1 {
                font-size: 1.8rem;
            }

            .charts-container {
                grid-template-columns: 1fr;
            }

            .stats-grid {
                grid-template-columns: repeat(2, 1fr);
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>📊 {{.RepoName}}</h1>
            <p>Code Analysis Visualization</p>
            <p style="font-size: 0.9rem; margin-top: 0.5rem;">Generated on {{.Analysis.AnalyzedAt}}</p>
        </header>

        <div class="content">
            <!-- Stats Overview -->
            <div class="stats-grid">
                <div class="stat-card">
                    <div class="stat-label">Total Files</div>
                    <div class="stat-value">{{.Analysis.Summary.TotalFiles}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Total Directories</div>
                    <div class="stat-value">{{.Analysis.Summary.TotalDirs}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Total Size</div>
                    <div class="stat-value">{{formatBytes .Analysis.Summary.TotalSize}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Languages</div>
                    <div class="stat-value">{{len .Analysis.Summary.Languages}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Max Depth</div>
                    <div class="stat-value">{{.Analysis.Summary.MaxDepth}}</div>
                </div>
                <div class="stat-card">
                    <div class="stat-label">Total Lines</div>
                    <div class="stat-value">{{.Analysis.Statistics.TotalLines}}</div>
                </div>
            </div>

            <!-- Charts -->
            <div class="section">
                <h2 class="section-title">📈 Language Distribution</h2>
                <div class="charts-container">
                    <div class="chart-box">
                        <div class="chart-title">Files by Language</div>
                        <canvas id="languagePieChart"></canvas>
                    </div>
                    <div class="chart-box">
                        <div class="chart-title">File Count Comparison</div>
                        <canvas id="languageBarChart"></canvas>
                    </div>
                </div>
            </div>

            <!-- Languages Detected -->
            {{if .Analysis.Summary.Languages}}
            <div class="section">
                <h2 class="section-title">💻 Languages Detected</h2>
                <div class="languages-badges">
                    {{range $lang, $count := .Analysis.Summary.Languages}}
                    <div class="language-badge">{{$lang}}: {{$count}} files</div>
                    {{end}}
                </div>
            </div>
            {{end}}

            <!-- Largest Files -->
            {{if .Analysis.Summary.LargestFiles}}
            <div class="section">
                <h2 class="section-title">📁 Largest Files</h2>
                <div class="file-list">
                    {{range .Analysis.Summary.LargestFiles}}
                    <div class="file-item">
                        <span class="file-name">{{relPath .Path}}</span>
                        <span class="file-size">{{formatBytes .Size}}</span>
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}

            <!-- Dependencies -->
            {{if .Analysis.Dependencies.PackageManagers}}
            <div class="section">
                <h2 class="section-title">📦 Dependencies</h2>
                {{range $name, $pm := .Analysis.Dependencies.PackageManagers}}
                <div class="dependencies-section">
                    <h3 class="dep-title">{{$pm.Name}}</h3>
                    <p style="color: #666; margin-bottom: 15px;">Config: {{range $pm.ConfigFiles}}{{relPath .}} {{end}}</p>
                    {{if $pm.Dependencies}}
                    <div class="dep-list">
                        {{range $dep, $ver := $pm.Dependencies}}
                        <div class="dep-item">
                            <span class="dep-name">{{$dep}}</span>
                            <span class="dep-version">{{$ver}}</span>
                        </div>
                        {{end}}
                    </div>
                    {{end}}
                </div>
                {{end}}
            </div>
            {{end}}

            <!-- Directory Structure -->
            <div class="section">
                <h2 class="section-title">🗂️ Directory Structure</h2>
                <div class="tree" id="directoryTree"></div>
            </div>
        </div>

        <footer>
            <p>Generated by Code Crawler &copy; 2025</p>
            <p style="margin-top: 0.5rem; font-size: 0.9rem;">Powered by Go & Chart.js</p>
        </footer>
    </div>

    <script>
        // Language distribution data
        const languageData = JSON.parse('{{toJSON .Analysis.Summary.Languages}}');
        const languageLabels = Object.keys(languageData);
        const languageCounts = Object.values(languageData);

        // Color palette
        const colors = [
            '#667eea', '#764ba2', '#f093fb', '#4facfe',
            '#43e97b', '#fa709a', '#fee140', '#30cfd0',
            '#a8edea', '#fed6e3', '#c1dfc4', '#deab6d'
        ];

        // Pie Chart
        const pieCtx = document.getElementById('languagePieChart').getContext('2d');
        new Chart(pieCtx, {
            type: 'pie',
            data: {
                labels: languageLabels,
                datasets: [{
                    data: languageCounts,
                    backgroundColor: colors,
                    borderWidth: 2,
                    borderColor: '#fff'
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: {
                        position: 'bottom',
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                const label = context.label || '';
                                const value = context.parsed || 0;
                                const total = context.dataset.data.reduce((a, b) => a + b, 0);
                                const percentage = ((value / total) * 100).toFixed(1);
                                return label + ': ' + value + ' files (' + percentage + '%)';
                            }
                        }
                    }
                }
            }
        });

        // Bar Chart
        const barCtx = document.getElementById('languageBarChart').getContext('2d');
        new Chart(barCtx, {
            type: 'bar',
            data: {
                labels: languageLabels,
                datasets: [{
                    label: 'Number of Files',
                    data: languageCounts,
                    backgroundColor: colors,
                    borderWidth: 0,
                    borderRadius: 8
                }]
            },
            options: {
                responsive: true,
                plugins: {
                    legend: {
                        display: false
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            precision: 0
                        }
                    }
                }
            }
        });

        // Directory Tree
        const treeData = JSON.parse('{{toJSON .Analysis.FileTree}}');
        const treeElement = document.getElementById('directoryTree');
        
        function renderNode(node, prefix = '', isLast = true) {
            if (!node) {
                return '<div class="tree-item">No files to display</div>';
            }

            let html = '';
            const connector = isLast ? '└── ' : '├── ';
            const extension = isLast ? '    ' : '│   ';
            
            if (prefix !== '') {
                if (node.is_dir) {
                    html += '<div class="tree-item tree-dir">' + prefix + connector + node.name + '/</div>';
                } else {
                    html += '<div class="tree-item tree-file">' + prefix + connector + node.name + '</div>';
                }
            } else {
                // Root node
                html += '<div class="tree-item tree-dir">' + node.name + '/</div>';
            }

            if (node.children && node.children.length > 0) {
                node.children.forEach((child, idx) => {
                    const isLastChild = idx === node.children.length - 1;
                    const newPrefix = prefix === '' ? '' : prefix + extension;
                    html += renderNode(child, newPrefix, isLastChild);
                });
            }

            return html;
        }

        if (treeData) {
            treeElement.innerHTML = renderNode(treeData);
        } else {
            treeElement.innerHTML = '<div class="tree-item">No files to display</div>';
        }
    </script>
</body>
</html>
`
