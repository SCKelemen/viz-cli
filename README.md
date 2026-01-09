# viz-cli

Command-line data visualization tool with dual-mode output: SVG for web and terminal for CLI.

## Features

- **Multiple Visualization Types**: Heatmaps, line graphs, bar charts, stat cards
- **Dual Output Modes**: SVG (vector graphics) and terminal (ASCII/Unicode with braille characters)
- **Enhanced Terminal Rendering**: Smooth braille character curves and ANSI color gradients
- **Interactive Dashboard**: Real-time TUI with bubbletea
- **Theme Support**: Default, midnight, nord, paper, wrapped themes
- **Data Input**: JSON files or stdin
- **Configurable**: Width, height, colors, and more

## Interactive Dashboard

### Simple Dashboard (Recommended)

The simple dashboard provides reliable visualization rendering with braille characters and ANSI colors:

```bash
go build simple_dashboard.go
./simple_dashboard
```

**Controls:**
- `t`: Toggle theme (blue/purple)
- `q`: Quit

**Features:**
- Contribution heatmap with color gradients
- Line graphs with smooth braille curves (⠀⠁⠂⠃⠄⠅⠆⠇)
- Colored bar charts
- Real-time updates
- Theme switching

## Installation

```bash
go install github.com/SCKelemen/viz-cli@latest
```

Or build from source:

```bash
git clone https://github.com/SCKelemen/viz-cli.git
cd viz-cli
go build
```

## Usage

```bash
viz-cli [options]

Options:
  -type string
        Visualization type: heatmap, line-graph, bar-chart, stat-card (default "heatmap")
  -format string
        Output format: svg, terminal (default "terminal")
  -data string
        Path to JSON data file (or use stdin with -)
  -theme string
        Theme: default, midnight, nord, paper, wrapped (default "default")
  -width int
        Width in pixels (SVG) or characters (terminal) (default 80)
  -height int
        Height in pixels (SVG) or characters (terminal) (default 24)
  -color string
        Primary color for visualization (hex format) (default "#3B82F6")
```

## Examples

### Terminal Output

```bash
# Heatmap from file
viz-cli -type heatmap -format terminal -data examples/heatmap.json

# Line graph from stdin
cat metrics.json | viz-cli -type line-graph -format terminal

# Bar chart with custom theme
viz-cli -type bar-chart -data repos.json -theme midnight
```

### SVG Output

```bash
# Generate SVG file
viz-cli -type heatmap -format svg -data examples/heatmap.json -width 400 -height 70 > heatmap.svg

# Pipe to SVG viewer
viz-cli -type line-graph -format svg -data metrics.json | display
```

## Data Formats

### Heatmap

```json
{
  "days": [
    {"date": "2024-01-01T00:00:00Z", "count": 10},
    {"date": "2024-01-02T00:00:00Z", "count": 15}
  ],
  "type": "linear"
}
```

### Line Graph

```json
{
  "points": [
    {"date": "2024-01-01T00:00:00Z", "value": 100},
    {"date": "2024-01-02T00:00:00Z", "value": 150}
  ],
  "color": "#3B82F6",
  "fillColor": "rgba(59, 130, 246, 0.1)",
  "useGradient": true
}
```

### Bar Chart

```json
{
  "bars": [
    {"value": 100, "secondary": 50, "label": "Item 1"},
    {"value": 80, "secondary": 30, "label": "Item 2"}
  ],
  "color": "#3B82F6",
  "stacked": true
}
```

### Stat Card

```json
{
  "title": "Total Commits",
  "value": "1,234",
  "subtitle": "past month",
  "color": "#3B82F6",
  "trendData": [
    {"date": "2024-01-01T00:00:00Z", "value": 10},
    {"date": "2024-01-02T00:00:00Z", "value": 15}
  ],
  "trendColor": "#3B82F6"
}
```

## Dependencies

- [github.com/SCKelemen/dataviz](https://github.com/SCKelemen/dataviz) - Visualization components
- [github.com/SCKelemen/design-system](https://github.com/SCKelemen/design-system) - Design tokens and themes
