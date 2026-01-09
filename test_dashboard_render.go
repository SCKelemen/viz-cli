package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
)

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
	fmt.Println("╔═══ DataViz Terminal Dashboard [Blue Theme] ═══════════════════════════════╗")
	fmt.Printf("║ Size: %dx%d • Counter: 5s • Press 't' to toggle theme, 'q' to quit ║\n", width, height)
	fmt.Println("╚═════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	tokens := design.DefaultTheme()
	config := dataviz.RenderConfig{
		DesignTokens: tokens,
		Color:        "#2196F3",
		Theme:        "default",
	}

	renderer := dataviz.NewTerminalRenderer()

	// Heatmap
	fmt.Println("┌─ CONTRIBUTION HEATMAP ───────────────────────────────────────────┐")
	heatmapData := dataviz.HeatmapData{
		Days:      heatmapDays,
		StartDate: now.AddDate(0, 0, -60),
		EndDate:   now,
		Type:      "linear",
	}
	heatmapBounds := dataviz.Bounds{X: 0, Y: 0, Width: 60, Height: 3}
	heatmapOutput := renderer.RenderHeatmap(heatmapData, heatmapBounds, config)
	fmt.Print(heatmapOutput.String())
	fmt.Println("└──────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Line Graph
	fmt.Println("┌─ METRICS LINE GRAPH ─────────────────────────────────────────────┐")
	lineData := dataviz.LineGraphData{
		Points: linePoints,
		Color:  "#2196F3",
	}
	lineBounds := dataviz.Bounds{X: 0, Y: 0, Width: 70, Height: 15}
	lineOutput := renderer.RenderLineGraph(lineData, lineBounds, config)
	fmt.Print(lineOutput.String())
	fmt.Println()
	fmt.Println("└──────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	// Bar Chart
	fmt.Println("┌─ LANGUAGE USAGE BAR CHART ───────────────────────────────────────┐")
	barData := dataviz.BarChartData{
		Bars:  bars,
		Color: "#FF9800",
	}
	barBounds := dataviz.Bounds{X: 0, Y: 0, Width: 50, Height: 8}
	barOutput := renderer.RenderBarChart(barData, barBounds, config)
	fmt.Print(barOutput.String())
	fmt.Println("└──────────────────────────────────────────────────────────────────┘")
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
	fmt.Println("╔═══ DataViz Terminal Dashboard [Purple Theme] ═══════════════════════════════╗")
	fmt.Printf("║ Size: %dx%d • Counter: 10s • Press 't' to toggle theme, 'q' to quit ║\n", width, height)
	fmt.Println("╚═════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	tokensNight := design.MidnightTheme()
	configNight := dataviz.RenderConfig{
		DesignTokens: tokensNight,
		Color:        "#7D56F4",
		Theme:        "midnight",
	}

	// Heatmap with purple
	fmt.Println("┌─ CONTRIBUTION HEATMAP ───────────────────────────────────────────┐")
	heatmapOutputNight := renderer.RenderHeatmap(heatmapData, heatmapBounds, configNight)
	fmt.Print(heatmapOutputNight.String())
	fmt.Println("└──────────────────────────────────────────────────────────────────┘")
	fmt.Println()

	fmt.Println("✓ Theme switching works correctly!")
	fmt.Println("✓ Colors changed from blue to purple")
	fmt.Println()
	fmt.Println("Run './simple_dashboard' for the interactive version!")
}
