package main

import (
	"fmt"

	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
)

func main() {
	// Test bar chart with various widths
	languages := []string{"Go", "TypeScript", "Python", "Rust", "JavaScript"}
	bars := make([]dataviz.BarData, len(languages))
	for i, lang := range languages {
		bars[i] = dataviz.BarData{
			Label: lang,
			Value: 100 - i*12,
		}
	}

	barData := dataviz.BarChartData{
		Bars:  bars,
		Color: "#FF9800",
	}

	tokens := design.DefaultTheme()
	config := dataviz.RenderConfig{
		DesignTokens: tokens,
		Color:        "#FF9800",
	}

	renderer := dataviz.NewTerminalRenderer()

	// Test with width 50
	fmt.Println("=== BAR CHART WIDTH 50 ===")
	bounds := dataviz.Bounds{X: 0, Y: 0, Width: 50, Height: 8}
	output := renderer.RenderBarChart(barData, bounds, config)
	fmt.Println(output.String())

	// Test with width 70
	fmt.Println("\n=== BAR CHART WIDTH 70 ===")
	bounds = dataviz.Bounds{X: 0, Y: 0, Width: 70, Height: 8}
	output = renderer.RenderBarChart(barData, bounds, config)
	fmt.Println(output.String())

	// Test with width 100
	fmt.Println("\n=== BAR CHART WIDTH 100 ===")
	bounds = dataviz.Bounds{X: 0, Y: 0, Width: 100, Height: 8}
	output = renderer.RenderBarChart(barData, bounds, config)
	fmt.Println(output.String())
}
