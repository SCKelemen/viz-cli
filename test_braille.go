package main

import (
	"fmt"

	"github.com/SCKelemen/dataviz"
)

func main() {
	// Test braille canvas directly
	canvas := dataviz.NewBrailleCanvas(40, 10)

	// Draw a simple diagonal line
	for i := 0; i < 40; i++ {
		x := i * 2
		y := i
		if y < 40 {
			canvas.SetPixel(x, y)
		}
	}

	output := canvas.Render()
	fmt.Println("=== BRAILLE CANVAS TEST ===")
	fmt.Println(output)
	fmt.Println("=== END BRAILLE TEST ===")

	// Test with points
	canvas2 := dataviz.NewBrailleCanvas(40, 10)
	points := []dataviz.Point{
		{X: 0, Y: 20},
		{X: 20, Y: 10},
		{X: 40, Y: 30},
		{X: 60, Y: 5},
	}
	canvas2.DrawCurve(points)

	output2 := canvas2.Render()
	fmt.Println("=== BRAILLE CURVE TEST ===")
	fmt.Println(output2)
	fmt.Println("=== END CURVE TEST ===")
}
