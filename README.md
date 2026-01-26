# viz-cli

> **⚠️ ARCHIVED: This repository has been merged into [dataviz](https://github.com/SCKelemen/dataviz)**
>
> Please use the new monorepo for all visualization needs. This repo is no longer maintained.
>
> **Migration:** Install from the new repo:
> ```bash
> go install github.com/SCKelemen/dataviz/cmd/viz-cli@latest
> ```
> Or build from source:
> ```bash
> git clone https://github.com/SCKelemen/dataviz
> cd dataviz
> go build -o viz-cli ./cmd/viz-cli
> ```
>
> See the [dataviz documentation](https://github.com/SCKelemen/dataviz) for updated usage instructions.

---

## Original README (Archived)

Interactive command-line data visualization tool with dual-mode output: SVG for web and terminal for CLI.

This tool uses the **[dataviz](https://github.com/SCKelemen/dataviz)** library for chart generation.

## DataViz Ecosystem

This tool has been consolidated into the dataviz monorepo. For current information, visit:

- **[dataviz](https://github.com/SCKelemen/dataviz)** - New monorepo with viz-cli, dataviz-mcp, and core library

### Migration Guide

**Old import/install:**
```bash
go install github.com/SCKelemen/viz-cli@latest
```

**New import/install:**
```bash
go install github.com/SCKelemen/dataviz/cmd/viz-cli@latest
```

All functionality from this repo is available in the new monorepo with the same CLI interface.

## Features (Archived)

- **Multiple Visualization Types**: Heatmaps, line graphs, bar charts, stat cards
- **Dual Output Modes**: SVG (vector graphics) and terminal (ASCII/Unicode with braille characters)
- **Enhanced Terminal Rendering**: Smooth braille character curves and ANSI color gradients
- **Interactive Dashboard**: Real-time TUI with bubbletea
- **Theme Support**: Default, midnight, nord, paper, wrapped themes
- **Data Input**: JSON files or stdin
- **Configurable**: Width, height, colors, and more

## License

BearWare 1.0 (MIT Compatible) 🐻

See [LICENSE](LICENSE) for the full license text.
