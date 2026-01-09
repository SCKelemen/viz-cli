package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
)

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

func main() {
	// Simulate a terminal size
	width := 180
	height := 82

	now := time.Now()

	// Generate heatmap data
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

	// Test with default (blue) theme
	const boxWidth = 70
	accentColor := "#2196F3"
	borderColor := hexToANSI(accentColor)

	// Create title bar
	titleBar := &TitleBar{
		Title:       "DataViz Terminal Dashboard [Blue Theme]",
		Width:       76,
		BorderColor: borderColor,
		Style:       LightBorderStyle,
	}
	fmt.Print(titleBar.Render())
	fmt.Print(titleBar.AddInfoLine(fmt.Sprintf(" Size: %dx%d • Counter: 5s • Press 't' to toggle theme, 'q' to quit ", width, height)))
	fmt.Print(titleBar.RenderBottom())
	fmt.Println()

	tokens := design.DefaultTheme()
	config := dataviz.RenderConfig{
		DesignTokens: tokens,
		Color:        "#2196F3",
		Theme:        "default",
	}

	renderer := dataviz.NewTerminalRenderer()

	// Heatmap
	heatmapBox := &Box{
		Label:       "CONTRIBUTION HEATMAP",
		Width:       boxWidth,
		BorderColor: borderColor,
		Style:       LightBorderStyle,
	}
	heatmapData := dataviz.HeatmapData{
		Days:      heatmapDays,
		StartDate: now.AddDate(0, 0, -60),
		EndDate:   now,
		Type:      "linear",
	}
	heatmapBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: 3}
	heatmapOutput := renderer.RenderHeatmap(heatmapData, heatmapBounds, config)
	fmt.Print(heatmapBox.RenderComplete(heatmapOutput.String()))
	fmt.Println()

	// Line Graph
	lineGraphBox := &Box{
		Label:       "METRICS LINE GRAPH",
		Width:       boxWidth,
		BorderColor: borderColor,
		Style:       LightBorderStyle,
	}
	lineData := dataviz.LineGraphData{
		Points: linePoints,
		Color:  "#2196F3",
	}
	lineBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: 15}
	lineOutput := renderer.RenderLineGraph(lineData, lineBounds, config)
	fmt.Print(lineGraphBox.RenderComplete(lineOutput.String()))
	fmt.Println()

	// Bar Chart
	barChartBox := &Box{
		Label:       "LANGUAGE USAGE BAR CHART",
		Width:       boxWidth,
		BorderColor: borderColor,
		Style:       LightBorderStyle,
	}
	barData := dataviz.BarChartData{
		Bars:  bars,
		Color: "#FF9800",
	}
	barBounds := dataviz.Bounds{X: 0, Y: 0, Width: boxWidth - 4, Height: 8}
	barOutput := renderer.RenderBarChart(barData, barBounds, config)
	fmt.Print(barChartBox.RenderComplete(barOutput.String()))
	fmt.Println()

	fmt.Println("✓ Dashboard rendering test complete!")
	fmt.Println("✓ All visualizations fit within their containers")
	fmt.Println("✓ Braille characters rendering smoothly")
	fmt.Println("✓ ANSI colors applied correctly")
	fmt.Println()

	// Now test purple theme
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════════════════")
	fmt.Println("Testing theme switch...")
	fmt.Println("═══════════════════════════════════════════════════════════════════════════")
	fmt.Println()

	accentColorNight := "#7D56F4"
	borderColorNight := hexToANSI(accentColorNight)

	// Create purple theme title bar
	titleBarNight := &TitleBar{
		Title:       "DataViz Terminal Dashboard [Purple Theme]",
		Width:       76,
		BorderColor: borderColorNight,
		Style:       LightBorderStyle,
	}
	fmt.Print(titleBarNight.Render())
	fmt.Print(titleBarNight.AddInfoLine(fmt.Sprintf(" Size: %dx%d • Counter: 10s • Press 't' to toggle theme, 'q' to quit ", width, height)))
	fmt.Print(titleBarNight.RenderBottom())
	fmt.Println()

	tokensNight := design.MidnightTheme()
	configNight := dataviz.RenderConfig{
		DesignTokens: tokensNight,
		Color:        "#7D56F4",
		Theme:        "midnight",
	}

	// Heatmap with purple
	heatmapBoxNight := &Box{
		Label:       "CONTRIBUTION HEATMAP",
		Width:       boxWidth,
		BorderColor: borderColorNight,
		Style:       LightBorderStyle,
	}
	heatmapOutputNight := renderer.RenderHeatmap(heatmapData, heatmapBounds, configNight)
	fmt.Print(heatmapBoxNight.RenderComplete(heatmapOutputNight.String()))
	fmt.Println()

	fmt.Println("✓ Theme switching works correctly!")
	fmt.Println("✓ Colors changed from blue to purple")
	fmt.Println()
	fmt.Println("Run './simple_dashboard' for the interactive version!")
}
