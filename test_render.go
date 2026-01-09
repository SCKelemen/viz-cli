package main

import (
	"fmt"
	"time"

	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
)

func main() {
	// Test heatmap
	now := time.Now()
	heatmapDays := make([]dataviz.ContributionDay, 30)
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -30+i)
		count := (i % 7) * 3
		heatmapDays[i] = dataviz.ContributionDay{
			Date:  date,
			Count: count,
		}
	}

	heatmapData := dataviz.HeatmapData{
		Days:      heatmapDays,
		StartDate: now.AddDate(0, 0, -30),
		EndDate:   now,
		Type:      "linear",
	}

	tokens := design.DefaultTheme()
	config := dataviz.RenderConfig{
		DesignTokens: tokens,
		Color:        "#40C463",
	}

	bounds := dataviz.Bounds{
		X:      0,
		Y:      0,
		Width:  80,
		Height: 10,
	}

	renderer := dataviz.NewTerminalRenderer()
	output := renderer.RenderHeatmap(heatmapData, bounds, config)

	fmt.Println("=== HEATMAP OUTPUT ===")
	fmt.Println(output.String())
	fmt.Println("=== END HEATMAP ===")

	// Test line graph
	linePoints := make([]dataviz.TimeSeriesData, 20)
	for i := 0; i < 20; i++ {
		linePoints[i] = dataviz.TimeSeriesData{
			Date:  now.AddDate(0, 0, -20+i),
			Value: 50 + i*2,
		}
	}

	lineData := dataviz.LineGraphData{
		Points: linePoints,
		Color:  "#2196F3",
	}

	lineBounds := dataviz.Bounds{
		X:      0,
		Y:      0,
		Width:  80,
		Height: 20,
	}

	lineOutput := renderer.RenderLineGraph(lineData, lineBounds, config)
	fmt.Println("=== LINE GRAPH OUTPUT ===")
	fmt.Println(lineOutput.String())
	fmt.Println("=== END LINE GRAPH ===")
}
