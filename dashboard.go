package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/SCKelemen/cli/renderer"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/dataviz"
	design "github.com/SCKelemen/design-system"
	"github.com/SCKelemen/layout"
)

type tickMsg time.Time

type viewMode int

const (
	viewHeatmap viewMode = iota
	viewLineGraph
	viewBarChart
	viewMulti
)

type dashboardModel struct {
	width      int
	height     int
	ready      bool
	counter    int
	mode       viewMode
	data       *dashboardData
	paused     bool
	colorTheme string
}

type dashboardData struct {
	heatmap   dataviz.HeatmapData
	lineGraph dataviz.LineGraphData
	barChart  dataviz.BarChartData
	lastUpdate time.Time
}

func initialDashboardModel() dashboardModel {
	return dashboardModel{
		mode:       viewMulti,
		data:       generateInitialData(),
		colorTheme: "default",
	}
}

func generateInitialData() *dashboardData {
	now := time.Now()

	// Generate heatmap data (last 52 weeks)
	heatmapDays := make([]dataviz.ContributionDay, 365)
	for i := 0; i < 365; i++ {
		date := now.AddDate(0, 0, -365+i)
		count := int(math.Abs(math.Sin(float64(i)/7)*20) + float64(rand.Intn(10)))
		heatmapDays[i] = dataviz.ContributionDay{
			Date:  date,
			Count: count,
		}
	}

	// Generate line graph data (last 30 days)
	linePoints := make([]dataviz.TimeSeriesData, 30)
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -30+i)
		value := int(50 + 30*math.Sin(float64(i)/5) + float64(rand.Intn(20)))
		linePoints[i] = dataviz.TimeSeriesData{
			Date:  date,
			Value: value,
		}
	}

	// Generate bar chart data
	languages := []string{"Go", "TypeScript", "Python", "Rust", "JavaScript"}
	bars := make([]dataviz.BarData, len(languages))
	for i, lang := range languages {
		bars[i] = dataviz.BarData{
			Label: lang,
			Value: 100 - i*15 + rand.Intn(20),
		}
	}

	return &dashboardData{
		heatmap: dataviz.HeatmapData{
			Days:      heatmapDays,
			StartDate: now.AddDate(0, 0, -365),
			EndDate:   now,
			Type:      "weeks",
		},
		lineGraph: dataviz.LineGraphData{
			Points:      linePoints,
			Color:       "#2196F3",
			UseGradient: true,
			Label:       "Metrics",
		},
		barChart: dataviz.BarChartData{
			Bars:  bars,
			Color: "#FF9800",
			Label: "Languages",
		},
		lastUpdate: now,
	}
}

func (m dashboardModel) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m dashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "1":
			m.mode = viewHeatmap
		case "2":
			m.mode = viewLineGraph
		case "3":
			m.mode = viewBarChart
		case "4", "m":
			m.mode = viewMulti
		case "p", " ":
			m.paused = !m.paused
		case "r":
			m.data = generateInitialData()
		case "t":
			// Toggle theme
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
		if !m.paused {
			m.counter++
			// Update data occasionally
			if m.counter%5 == 0 {
				m.updateData()
			}
		}
		return m, tickCmd()
	}

	return m, nil
}

func (m *dashboardModel) updateData() {
	now := time.Now()

	// Add new point to line graph (shift old data)
	if len(m.data.lineGraph.Points) > 0 {
		m.data.lineGraph.Points = append(m.data.lineGraph.Points[1:], dataviz.TimeSeriesData{
			Date:  now,
			Value: int(50 + 30*math.Sin(float64(m.counter)/5) + float64(rand.Intn(20))),
		})
	}

	// Update bar chart values slightly
	for i := range m.data.barChart.Bars {
		change := rand.Intn(11) - 5 // -5 to +5
		m.data.barChart.Bars[i].Value += change
		if m.data.barChart.Bars[i].Value < 10 {
			m.data.barChart.Bars[i].Value = 10
		}
		if m.data.barChart.Bars[i].Value > 200 {
			m.data.barChart.Bars[i].Value = 200
		}
	}

	m.data.lastUpdate = now
}

func (m dashboardModel) View() string {
	if !m.ready {
		return "Initializing dashboard...\n\nPress any key to continue"
	}

	screen := renderer.NewScreen(m.width, m.height)
	ctx := layout.NewLayoutContext(float64(m.width), float64(m.height), 16)

	// Get theme with distinct colors
	var tokens *design.DesignTokens
	var accentColor string
	if m.colorTheme == "midnight" {
		tokens = design.MidnightTheme()
		accentColor = "#7D56F4" // Purple for midnight theme
	} else {
		tokens = design.DefaultTheme()
		accentColor = "#2196F3" // Blue for default theme
	}

	config := dataviz.RenderConfig{
		DesignTokens: tokens,
		Color:        accentColor,
		Theme:        m.colorTheme,
	}

	switch m.mode {
	case viewHeatmap:
		return m.renderSingleView(screen, ctx, config, "Contribution Heatmap", m.renderHeatmap)
	case viewLineGraph:
		return m.renderSingleView(screen, ctx, config, "Metrics Over Time", m.renderLineGraph)
	case viewBarChart:
		return m.renderSingleView(screen, ctx, config, "Language Usage", m.renderBarChart)
	case viewMulti:
		return m.renderMultiView(screen, ctx, config)
	}

	return ""
}

func (m dashboardModel) renderSingleView(screen *renderer.Screen, ctx *layout.LayoutContext, config dataviz.RenderConfig, title string, renderFunc func(dataviz.Bounds, dataviz.RenderConfig) string) string {
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Vw(100),
			Height:        layout.Vh(100),
			Padding:       layout.Uniform(layout.Ch(1)),
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header
	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Height:  layout.Ch(3),
		},
	}
	white, _ := color.ParseColor("#FAFAFA")
	accent, _ := color.ParseColor(config.Color)
	headerStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &accent,
	}
	headerStyle.WithBorder(renderer.RoundedBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = fmt.Sprintf(" %s • %dx%d • %s", title, m.width, m.height, m.getStatusText())
	rootStyled.AddChild(headerStyled)

	// Visualization
	vizNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Margin:  layout.Spacing{Top: layout.Ch(1)},
		},
	}
	vizStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &accent,
	}
	vizStyle.WithBorder(renderer.RoundedBorder)
	vizStyled := renderer.NewStyledNode(vizNode, vizStyle)

	// Render visualization
	vizBounds := dataviz.Bounds{
		X:      0,
		Y:      0,
		Width:  m.width - 4,
		Height: m.height - 8,
	}
	vizContent := renderFunc(vizBounds, config)
	vizStyled.Content = "\n" + vizContent
	rootStyled.AddChild(vizStyled)

	// Controls footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Height:  layout.Ch(2),
			Margin:  layout.Spacing{Top: layout.Ch(1)},
		},
	}
	footerStyled := renderer.NewStyledNode(footerNode, nil)
	footerStyled.Content = m.getControlsText()
	rootStyled.AddChild(footerStyled)

	constraints := layout.Tight(float64(m.width), float64(m.height))
	layout.Layout(root, constraints, ctx)
	screen.Render(rootStyled)

	return screen.String()
}

func (m dashboardModel) renderMultiView(screen *renderer.Screen, ctx *layout.LayoutContext, config dataviz.RenderConfig) string {
	root := &layout.Node{
		Style: layout.Style{
			Display:             layout.DisplayGrid,
			Width:               layout.Vw(100),
			Height:              layout.Vh(100),
			GridTemplateColumns: []layout.GridTrack{layout.FractionTrack(1)},
			GridTemplateRows:    []layout.GridTrack{layout.FixedTrack(layout.Ch(3)), layout.FractionTrack(1), layout.FractionTrack(1), layout.FractionTrack(1), layout.FixedTrack(layout.Ch(2))},
			GridGap:             layout.Ch(1),
			Padding:             layout.Uniform(layout.Ch(1)),
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	white, _ := color.ParseColor("#FAFAFA")
	accent, _ := color.ParseColor(config.Color)

	// Header
	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
		},
	}
	headerStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &accent,
	}
	headerStyle.WithBorder(renderer.RoundedBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = fmt.Sprintf(" DataViz Dashboard • %dx%d • %s", m.width, m.height, m.getStatusText())
	rootStyled.AddChild(headerStyled)

	// Heatmap panel
	heatmapNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
		},
	}
	blue, _ := color.ParseColor("#2196F3")
	heatmapStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &blue,
	}
	heatmapStyle.WithBorder(renderer.RoundedBorder)
	heatmapStyled := renderer.NewStyledNode(heatmapNode, heatmapStyle)
	heatmapBounds := dataviz.Bounds{X: 0, Y: 0, Width: m.width - 4, Height: 10}
	heatmapContent := m.renderHeatmap(heatmapBounds, config)
	heatmapStyled.Content = " Contributions\n" + heatmapContent
	rootStyled.AddChild(heatmapStyled)

	// Line graph panel
	lineNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
		},
	}
	green, _ := color.ParseColor("#4CAF50")
	lineStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &green,
	}
	lineStyle.WithBorder(renderer.RoundedBorder)
	lineStyled := renderer.NewStyledNode(lineNode, lineStyle)
	lineBounds := dataviz.Bounds{X: 0, Y: 0, Width: m.width - 4, Height: (m.height - 10) / 3}
	lineContent := m.renderLineGraph(lineBounds, config)
	lineStyled.Content = " Metrics\n" + lineContent
	rootStyled.AddChild(lineStyled)

	// Bar chart panel
	barNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
		},
	}
	orange, _ := color.ParseColor("#FF9800")
	barStyle := &renderer.Style{
		Foreground:  &white,
		BorderColor: &orange,
	}
	barStyle.WithBorder(renderer.RoundedBorder)
	barStyled := renderer.NewStyledNode(barNode, barStyle)
	barBounds := dataviz.Bounds{X: 0, Y: 0, Width: m.width - 4, Height: (m.height - 10) / 3}
	barContent := m.renderBarChart(barBounds, config)
	barStyled.Content = " Languages\n" + barContent
	rootStyled.AddChild(barStyled)

	// Controls footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
		},
	}
	footerStyled := renderer.NewStyledNode(footerNode, nil)
	footerStyled.Content = m.getControlsText()
	rootStyled.AddChild(footerStyled)

	constraints := layout.Tight(float64(m.width), float64(m.height))
	layout.Layout(root, constraints, ctx)
	screen.Render(rootStyled)

	return screen.String()
}

func (m dashboardModel) renderHeatmap(bounds dataviz.Bounds, config dataviz.RenderConfig) string {
	renderer := dataviz.NewTerminalRenderer()
	output := renderer.RenderHeatmap(m.data.heatmap, bounds, config)
	return output.String()
}

func (m dashboardModel) renderLineGraph(bounds dataviz.Bounds, config dataviz.RenderConfig) string {
	renderer := dataviz.NewTerminalRenderer()
	output := renderer.RenderLineGraph(m.data.lineGraph, bounds, config)
	return output.String()
}

func (m dashboardModel) renderBarChart(bounds dataviz.Bounds, config dataviz.RenderConfig) string {
	renderer := dataviz.NewTerminalRenderer()
	output := renderer.RenderBarChart(m.data.barChart, bounds, config)
	return output.String()
}

func (m dashboardModel) getStatusText() string {
	status := "Running"
	if m.paused {
		status = "Paused"
	}
	return fmt.Sprintf("%s • %s theme • %ds", status, m.colorTheme, m.counter)
}

func (m dashboardModel) getControlsText() string {
	return " 1:Heatmap 2:LineGraph 3:BarChart 4:Multi • p:Pause r:Refresh t:Theme q:Quit"
}

func main() {
	p := tea.NewProgram(initialDashboardModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
