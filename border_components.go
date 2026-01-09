package main

import (
	"strings"

	"github.com/SCKelemen/text"
)

// BorderStyle defines the characters used for box drawing
type BorderStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	// Title bar style
	TitleTopLeft     string
	TitleTopRight    string
	TitleBottomLeft  string
	TitleBottomRight string
	TitleHorizontal  string
	TitleVertical    string
}

var (
	// LightBorderStyle uses light box-drawing characters
	LightBorderStyle = BorderStyle{
		TopLeft:          "┌",
		TopRight:         "┐",
		BottomLeft:       "└",
		BottomRight:      "┘",
		Horizontal:       "─",
		Vertical:         "│",
		TitleTopLeft:     "╔",
		TitleTopRight:    "╗",
		TitleBottomLeft:  "╚",
		TitleBottomRight: "╝",
		TitleHorizontal:  "═",
		TitleVertical:    "║",
	}
)

// TitleBar creates a title bar with borders
type TitleBar struct {
	Title       string
	Width       int
	BorderColor string
	Style       BorderStyle
}

// Render creates the title bar output
func (tb *TitleBar) Render() string {
	var output strings.Builder

	colorCode := tb.BorderColor
	resetCode := "\x1b[0m"

	// Top line: ╔═══ Title ═════════════════════════════╗
	titlePrefix := tb.Style.TitleTopLeft + strings.Repeat(tb.Style.TitleHorizontal, 3) + " "
	titleSuffix := tb.Style.TitleTopRight
	prefixLen := 5 // "╔═══ "
	suffixLen := 1 // "╗"

	titlePadding := tb.Width - len(tb.Title) - prefixLen - suffixLen
	output.WriteString(colorCode + titlePrefix + resetCode + tb.Title + colorCode)
	if titlePadding > 0 {
		output.WriteString(strings.Repeat(tb.Style.TitleHorizontal, titlePadding))
	}
	output.WriteString(titleSuffix + resetCode + "\n")

	return output.String()
}

// AddInfoLine adds a content line between borders
func (tb *TitleBar) AddInfoLine(content string) string {
	colorCode := tb.BorderColor
	resetCode := "\x1b[0m"

	return colorCode + tb.Style.TitleVertical + resetCode + content + colorCode + tb.Style.TitleVertical + resetCode + "\n"
}

// RenderBottom creates the bottom border of the title bar
func (tb *TitleBar) RenderBottom() string {
	colorCode := tb.BorderColor
	resetCode := "\x1b[0m"

	return colorCode + tb.Style.TitleBottomLeft + strings.Repeat(tb.Style.TitleHorizontal, tb.Width-2) + tb.Style.TitleBottomRight + resetCode + "\n"
}

// Box creates a bordered box with optional title label
type Box struct {
	Label       string
	Width       int
	BorderColor string
	Style       BorderStyle
}

// RenderTop creates the top border with optional label
func (b *Box) RenderTop() string {
	colorCode := b.BorderColor
	resetCode := "\x1b[0m"

	if b.Label == "" {
		// Simple top border without label
		return colorCode + b.Style.TopLeft + strings.Repeat(b.Style.Horizontal, b.Width-2) + b.Style.TopRight + resetCode + "\n"
	}

	// Top border with label: ┌─ LABEL ──────────────┐
	labelPrefix := b.Style.TopLeft + b.Style.Horizontal + " "
	labelSuffix := b.Style.TopRight

	prefixLen := 3 // "┌─ "
	suffixLen := 1 // "┐"

	hLineLen := b.Width - len(b.Label) - prefixLen - suffixLen
	return colorCode + labelPrefix + resetCode + b.Label + colorCode + strings.Repeat(b.Style.Horizontal, hLineLen) + labelSuffix + resetCode + "\n"
}

// RenderBottom creates the bottom border
func (b *Box) RenderBottom() string {
	colorCode := b.BorderColor
	resetCode := "\x1b[0m"

	return colorCode + b.Style.BottomLeft + strings.Repeat(b.Style.Horizontal, b.Width-2) + b.Style.BottomRight + resetCode + "\n"
}

// WrapContent wraps content lines with left and right borders
func (b *Box) WrapContent(content string) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder

	colorCode := b.BorderColor
	resetCode := "\x1b[0m"

	// Content width is box width minus borders (│ on each side) and padding (1 space on each side)
	contentWidth := b.Width - 4

	// Create text handler for proper Unicode width measurement
	txt := text.NewTerminal()

	for _, line := range lines {
		// Skip completely empty lines
		if line == "" {
			continue
		}

		// Remove ANSI color codes and measure actual display width
		stripped := stripANSI(line)
		displayWidth := int(txt.Width(stripped))

		// Truncate if line is too long
		if displayWidth > contentWidth {
			line = truncateANSI(line, contentWidth)
			displayWidth = contentWidth
		}

		// Add left border and padding with color
		if colorCode != "" {
			result.WriteString(colorCode)
		}
		result.WriteString(b.Style.Vertical)
		if colorCode != "" {
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
		if colorCode != "" {
			result.WriteString(colorCode)
		}
		result.WriteString(b.Style.Vertical)
		if colorCode != "" {
			result.WriteString(resetCode)
		}
		result.WriteString("\n")
	}

	return result.String()
}

// RenderComplete renders a complete box with top, content, and bottom
func (b *Box) RenderComplete(content string) string {
	var output strings.Builder
	output.WriteString(b.RenderTop())
	output.WriteString(b.WrapContent(content))
	output.WriteString(b.RenderBottom())
	return output.String()
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
