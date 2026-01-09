package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
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

	// Get theme
	var accentColor string
	var tokens *design.DesignTokens
	var themeIndicator string
	if m.colorTheme == "midnight" {
		tokens = design.MidnightTheme()
		accentColor = "#7D56F4" // Purple
		themeIndicator = "[Purple Theme]"
	} else {
		tokens = design.DefaultTheme()
		accentColor = "#2196F3" // Blue
		themeIndicator = "[Blue Theme]"
	}

	config := dataviz.RenderConfig{
		DesignTokens: tokens,
		Color:        accentColor,
		Theme:        m.colorTheme,
	}

	renderer := dataviz.NewTerminalRenderer()

	// Convert accent color to ANSI code for borders
	borderColor := hexToANSI(accentColor)

	// Create title bar
	titleText := fmt.Sprintf("DataViz Terminal Dashboard %s", themeIndicator)
	titleBar := &TitleBar{
		Title:       titleText,
		Width:       76,
		BorderColor: borderColor,
		Style:       LightBorderStyle,
	}
	output += titleBar.Render()
	output += titleBar.AddInfoLine(fmt.Sprintf(" Size: %dx%d • Counter: %ds • Press 't' to toggle theme, 'q' to quit ", m.width, m.height, m.counter))
	output += titleBar.RenderBottom()
	output += "\n"

	// Heatmap
	heatmapBox := &Box{
		Label:       "CONTRIBUTION HEATMAP",
		Width:       boxWidth,
		BorderColor: borderColor,
		Style:       LightBorderStyle,
	}
	heatmapBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: 3}
	heatmapOutput := renderer.RenderHeatmap(m.heatmap, heatmapBounds, config)
	output += heatmapBox.RenderComplete(heatmapOutput.String())
	output += "\n"

	// Line Graph
	lineGraphBox := &Box{
		Label:       "METRICS LINE GRAPH",
		Width:       boxWidth,
		BorderColor: borderColor,
		Style:       LightBorderStyle,
	}
	lineHeight := 15
	if m.height > 40 {
		lineHeight = 20
	}
	lineBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: lineHeight}
	lineOutput := renderer.RenderLineGraph(m.lineGraph, lineBounds, config)
	output += lineGraphBox.RenderComplete(lineOutput.String())
	output += "\n"

	// Bar Chart
	barChartBox := &Box{
		Label:       "LANGUAGE USAGE BAR CHART",
		Width:       boxWidth,
		BorderColor: borderColor,
		Style:       LightBorderStyle,
	}
	barBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: 8}
	barOutput := renderer.RenderBarChart(m.barChart, barBounds, config)
	output += barChartBox.RenderComplete(barOutput.String())

	return output
}

func main() {
	p := tea.NewProgram(initialSimpleModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
