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
	"github.com/SCKelemen/text"
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

// wrapInBorders adds left and right borders to content with optional color
func wrapInBorders(content string, boxWidth int, borderColor string) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder

	// Content width is box width minus borders (│ on each side) and padding (1 space on each side)
	contentWidth := boxWidth - 4

	// ANSI color codes for borders
	colorCode := ""
	resetCode := "\x1b[0m"
	if borderColor != "" {
		colorCode = borderColor
	}

	// Create text handler for proper Unicode width measurement
	txt := text.NewTerminal()

	for _, line := range lines {
		// Skip completely empty lines
		if line == "" {
			continue
		}

		// Remove ANSI color codes and measure actual display width using text package
		stripped := stripANSI(line)
		displayWidth := int(txt.Width(stripped))

		// Truncate if line is too long
		if displayWidth > contentWidth {
			// Need to truncate while preserving ANSI codes
			line = truncateANSI(line, contentWidth)
			displayWidth = contentWidth
		}

		// Add left border and padding with color
		if borderColor != "" {
			result.WriteString(colorCode)
		}
		result.WriteString("│")
		if borderColor != "" {
			result.WriteString(resetCode)
		}
		result.WriteString(" ")
		result.WriteString(line)

		// Ensure any open ANSI codes are closed before padding
		if hasUnclosedANSI(line) {
			result.WriteString(resetCode)
		}

		// Add right padding - always pad to exact width for alignment
		paddingNeeded := contentWidth - displayWidth
		if paddingNeeded > 0 {
			result.WriteString(strings.Repeat(" ", paddingNeeded))
		}

		result.WriteString(" ")
		if borderColor != "" {
			result.WriteString(colorCode)
		}
		result.WriteString("│")
		if borderColor != "" {
			result.WriteString(resetCode)
		}
		result.WriteString("\n")
	}

	return result.String()
}

// hasUnclosedANSI checks if a string has unclosed ANSI codes
func hasUnclosedANSI(s string) bool {
	openCount := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			openCount++
		} else if openCount > 0 && s[i] == 'm' {
			openCount--
		}
	}
	return openCount > 0
}

// truncateANSI truncates a string to a certain visual width while preserving ANSI codes
func truncateANSI(s string, maxWidth int) string {
	var result strings.Builder
	visualWidth := 0
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
			result.WriteByte(s[i])
		} else if inEscape {
			result.WriteByte(s[i])
			if s[i] == 'm' {
				inEscape = false
			}
		} else {
			if visualWidth >= maxWidth {
				break
			}
			result.WriteByte(s[i])
			visualWidth++
		}
	}

	// Close any open ANSI codes
	if inEscape || hasUnclosedANSI(result.String()) {
		result.WriteString("\x1b[0m")
	}

	return result.String()
}

// stripANSI removes ANSI escape codes for measuring display width
func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false

	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
		} else if inEscape {
			if r == 'm' {
				inEscape = false
			}
			// Skip all characters while in escape sequence
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// hexToANSI converts a hex color to ANSI escape code
func hexToANSI(hexColor string) string {
	if hexColor == "" {
		return ""
	}
	// Remove # if present
	if len(hexColor) > 0 && hexColor[0] == '#' {
		hexColor = hexColor[1:]
	}
	if len(hexColor) != 6 {
		return ""
	}

	// Parse RGB values
	var r, g, b int
	fmt.Sscanf(hexColor, "%02x%02x%02x", &r, &g, &b)

	// Return ANSI TrueColor escape code
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
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

	// Convert accent color to ANSI code for borders
	borderColor := hexToANSI(accentColor)

	// Heatmap
	output += "┌─ CONTRIBUTION HEATMAP ───────────────────────────────────────────┐\n"
	heatmapBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: 3}
	heatmapOutput := renderer.RenderHeatmap(m.heatmap, heatmapBounds, config)
	output += wrapInBorders(heatmapOutput.String(), boxWidth, borderColor)
	output += "└──────────────────────────────────────────────────────────────────┘\n\n"

	// Line Graph
	output += "┌─ METRICS LINE GRAPH ─────────────────────────────────────────────┐\n"
	lineHeight := 15
	if m.height > 40 {
		lineHeight = 20
	}
	lineBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: lineHeight}
	lineOutput := renderer.RenderLineGraph(m.lineGraph, lineBounds, config)
	output += wrapInBorders(lineOutput.String(), boxWidth, borderColor)
	output += "└──────────────────────────────────────────────────────────────────┘\n\n"

	// Bar Chart
	output += "┌─ LANGUAGE USAGE BAR CHART ───────────────────────────────────────┐\n"
	barBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: 8}
	barOutput := renderer.RenderBarChart(m.barChart, barBounds, config)
	output += wrapInBorders(barOutput.String(), boxWidth, borderColor)
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
