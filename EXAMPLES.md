# Example Usage Scenarios

## Scenario 1: New Team Member Onboarding

You're joining a new team and need to understand the codebase quickly:

```bash
# Analyze the repository
cd ~/projects/team-repo
code-crawler -path . -output ~/onboarding-analysis

# Open the visualization
open onboarding-analysis/visualization.html
```

**What you'll learn:**
- Project structure and organization
- Main programming languages used
- Key dependencies and their versions
- Largest files (often the most complex)
- Deepest directory paths

## Scenario 2: Code Audit Before Refactoring

Planning a major refactoring and need to assess the current state:

```bash
# Generate baseline analysis
code-crawler -path ~/projects/legacy-app -output ./before-refactor

# After refactoring
code-crawler -path ~/projects/legacy-app -output ./after-refactor

# Compare the analysis.json files
```

## Scenario 3: Multi-Repository Analysis

Analyzing multiple microservices:

```bash
#!/bin/bash
# analyze-all.sh

for repo in service-a service-b service-c; do
  echo "Analyzing $repo..."
  code-crawler -path ~/projects/$repo -output ./analysis/$repo
done

echo "All analyses complete!"
```

## Scenario 4: CI/CD Integration

Track code metrics over time:

```bash
# In your CI pipeline (.github/workflows/analysis.yml)
name: Code Analysis
on: [push]
jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Download Code Crawler
        run: |
          curl -L https://github.com/shalinnaidoo/code-crawler/releases/latest/download/code-crawler-linux -o code-crawler
          chmod +x code-crawler
      - name: Analyze Code
        run: ./code-crawler -path . -output ./analysis
      - name: Upload Analysis
        uses: actions/upload-artifact@v2
        with:
          name: code-analysis
          path: analysis/
```

## Scenario 5: Documentation Generation

Generate up-to-date architecture documentation:

```bash
# Run analysis
code-crawler -path . -output ./docs/analysis

# The visualization.html can be included in your docs site
# The analysis.json can be used to generate custom reports
```

## Scenario 6: Dependency Audit

Check what dependencies your project uses:

```bash
code-crawler -path ~/projects/my-app

# Then check analysis.json for the dependencies section
cat code-analysis/analysis.json | jq '.dependencies'
```

## Scenario 7: Language Distribution Report

See which languages dominate your codebase:

```bash
code-crawler -path . -output ./metrics

# Extract language stats
cat ./metrics/analysis.json | jq '.summary.languages'
```

## Scenario 8: Finding Large Files

Identify files that might need splitting:

```bash
code-crawler -path .

# Check the visualization for "Largest Files" section
# Or query the JSON:
cat code-analysis/analysis.json | jq '.summary.largest_files'
```

## Scenario 9: Custom Analysis Script

Use the JSON output for custom reporting:

```python
# analyze.py
import json

with open('code-analysis/analysis.json') as f:
    data = json.load(f)

print(f"Total Files: {data['summary']['total_files']}")
print(f"Total Lines: {data['statistics']['total_lines']}")
print(f"Languages: {', '.join(data['summary']['languages'].keys())}")

# Calculate complexity score
complexity = (
    data['summary']['max_depth'] * 10 +
    data['summary']['total_files'] +
    len(data['dependencies']['external_deps'])
)
print(f"Complexity Score: {complexity}")
```

```bash
code-crawler -path ~/project
python analyze.py
```

## Scenario 10: Pre-Deployment Check

Verify project state before deployment:

```bash
#!/bin/bash
# pre-deploy.sh

code-crawler -path . -output ./pre-deploy-check

# Check if there are too many dependencies
DEPS=$(cat pre-deploy-check/analysis.json | jq '.dependencies.external_deps | length')
if [ $DEPS -gt 100 ]; then
  echo "Warning: High number of dependencies ($DEPS)"
fi

# Check total lines of code
LINES=$(cat pre-deploy-check/analysis.json | jq '.statistics.total_lines')
echo "Total lines of code: $LINES"
```
