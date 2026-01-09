package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
)

type tickMsg time.Time

type simpleModel struct {
	width      int
	height     int
	ready      bool
	counter    int
	colorTheme string
	heatmap    dataviz.HeatmapData
	lineGraph  dataviz.LineGraphData
	barChart   dataviz.BarChartData
}

func initialSimpleModel() simpleModel {
	now := time.Now()

	// Generate simple heatmap data
	heatmapDays := make([]dataviz.ContributionDay, 60)
	for i := 0; i < 60; i++ {
		date := now.AddDate(0, 0, -60+i)
		count := int(math.Abs(math.Sin(float64(i)/7)*20) + float64(rand.Intn(10)))
		heatmapDays[i] = dataviz.ContributionDay{
			Date:  date,
			Count: count,
		}
	}

	// Generate line graph data
	linePoints := make([]dataviz.TimeSeriesData, 30)
	for i := 0; i < 30; i++ {
		linePoints[i] = dataviz.TimeSeriesData{
			Date:  now.AddDate(0, 0, -30+i),
			Value: 30 + int(20*math.Sin(float64(i)/3)) + rand.Intn(10),
		}
	}

	// Generate bar chart data
	languages := []string{"Go", "TypeScript", "Python", "Rust", "JavaScript"}
	bars := make([]dataviz.BarData, len(languages))
	for i, lang := range languages {
		bars[i] = dataviz.BarData{
			Label: lang,
			Value: 100 - i*12,
		}
	}

	return simpleModel{
		colorTheme: "default",
		heatmap: dataviz.HeatmapData{
			Days:      heatmapDays,
			StartDate: now.AddDate(0, 0, -60),
			EndDate:   now,
			Type:      "linear",
		},
		lineGraph: dataviz.LineGraphData{
			Points: linePoints,
			Color:  "#2196F3",
		},
		barChart: dataviz.BarChartData{
			Bars:  bars,
			Color: "#FF9800",
		},
	}
}

func (m simpleModel) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m simpleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "t":
			if m.colorTheme == "default" {
				m.colorTheme = "midnight"
			} else {
				m.colorTheme = "default"
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tickMsg:
		m.counter++
		return m, tickCmd()
	}

	return m, nil
}

// wrapInBorders adds left and right borders to content
func wrapInBorders(content string, boxWidth int) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder

	// Content width is box width minus borders (│ on each side) and padding (1 space on each side)
	contentWidth := boxWidth - 4

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Remove ANSI color codes to measure actual display width
		displayWidth := len(stripANSI(line))

		// Add left border and padding
		result.WriteString("│ ")
		result.WriteString(line)

		// Add right padding and border
		paddingNeeded := contentWidth - displayWidth
		if paddingNeeded > 0 {
			result.WriteString(strings.Repeat(" ", paddingNeeded))
		}
		result.WriteString(" │\n")
	}

	return result.String()
}

// stripANSI removes ANSI escape codes for measuring display width
func stripANSI(s string) string {
	// Simple ANSI code stripper - matches \x1b[...m patterns
	result := ""
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
		} else if inEscape && s[i] == 'm' {
			inEscape = false
		} else if !inEscape {
			result += string(s[i])
		}
	}
	return result
}

func (m simpleModel) View() string {
	if !m.ready {
		return "Initializing...\n"
	}

	var output string

	// Box width for all containers
	const boxWidth = 70

	// Title
	themeIndicator := "[Blue Theme]"
	if m.colorTheme == "midnight" {
		themeIndicator = "[Purple Theme]"
	}
	output += fmt.Sprintf("╔═══ DataViz Terminal Dashboard %s ═══════════════════════════════╗\n", themeIndicator)
	output += fmt.Sprintf("║ Size: %dx%d • Counter: %ds • Press 't' to toggle theme, 'q' to quit ║\n", m.width, m.height, m.counter)
	output += "╚═════════════════════════════════════════════════════════════════════════╝\n\n"

	// Get theme
	var accentColor string
	var tokens *design.DesignTokens
	if m.colorTheme == "midnight" {
		tokens = design.MidnightTheme()
		accentColor = "#7D56F4" // Purple
	} else {
		tokens = design.DefaultTheme()
		accentColor = "#2196F3" // Blue
	}

	config := dataviz.RenderConfig{
		DesignTokens: tokens,
		Color:        accentColor,
		Theme:        m.colorTheme,
	}

	renderer := dataviz.NewTerminalRenderer()

	// Heatmap
	output += "┌─ CONTRIBUTION HEATMAP ───────────────────────────────────────────┐\n"
	heatmapBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: 3}
	heatmapOutput := renderer.RenderHeatmap(m.heatmap, heatmapBounds, config)
	output += wrapInBorders(heatmapOutput.String(), boxWidth)
	output += "└──────────────────────────────────────────────────────────────────┘\n\n"

	// Line Graph
	output += "┌─ METRICS LINE GRAPH ─────────────────────────────────────────────┐\n"
	lineHeight := 15
	if m.height > 40 {
		lineHeight = 20
	}
	lineBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: lineHeight}
	lineOutput := renderer.RenderLineGraph(m.lineGraph, lineBounds, config)
	output += wrapInBorders(lineOutput.String(), boxWidth)
	output += "└──────────────────────────────────────────────────────────────────┘\n\n"

	// Bar Chart
	output += "┌─ LANGUAGE USAGE BAR CHART ───────────────────────────────────────┐\n"
	barBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: 8}
	barOutput := renderer.RenderBarChart(m.barChart, barBounds, config)
	output += wrapInBorders(barOutput.String(), boxWidth)
	output += "└──────────────────────────────────────────────────────────────────┘\n"

	return output
}

func main() {
	p := tea.NewProgram(initialSimpleModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
