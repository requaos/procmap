package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/req/procmap/proc"
)

// Bubble represents a positioned bubble on the screen
type Bubble struct {
	X         float64
	Y         float64
	Radius    float64
	Process   proc.ProcessInfo
	Value     float64
	ColorHex  string
}

// getValue returns the metric value based on sort mode
func getValue(p proc.ProcessInfo, sortMode string) float64 {
	switch sortMode {
	case "cpu":
		return p.CPUPercent
	case "memory":
		return float64(p.MemPercent)
	case "threads":
		return float64(p.NumThreads)
	default:
		return p.CPUPercent
	}
}

// calculateBubbleSize maps a metric value to a radius
func calculateBubbleSize(value, maxValue, minValue float64, screenWidth, screenHeight int) float64 {
	if maxValue == minValue {
		return 3.0
	}

	// Normalize value to 0-1 range
	normalized := (value - minValue) / (maxValue - minValue)

	// Calculate max radius based on screen size
	maxRadius := math.Min(float64(screenWidth)/8.0, float64(screenHeight)/6.0)
	minRadius := 2.0

	// Map to radius range with sqrt for better visual distribution
	radius := minRadius + math.Sqrt(normalized)*(maxRadius-minRadius)
	return radius
}

// getColor returns a color based on the percentile of the value
func getColor(value, maxValue, minValue float64) string {
	if maxValue == minValue {
		return "12" // cyan
	}

	normalized := (value - minValue) / (maxValue - minValue)

	if normalized > 0.75 {
		return "9" // red
	} else if normalized > 0.5 {
		return "208" // orange
	} else if normalized > 0.25 {
		return "11" // yellow
	} else {
		return "10" // green
	}
}

// packBubbles implements a simple circle packing algorithm
func packBubbles(processes []proc.ProcessInfo, sortMode string, maxBubbles, width, height int) []Bubble {
	if len(processes) == 0 {
		return []Bubble{}
	}

	// Take top N processes
	count := maxBubbles
	if count > len(processes) {
		count = len(processes)
	}
	topProcs := processes[:count]

	// Get values and find min/max for sizing
	values := make([]float64, count)
	maxValue := 0.0
	minValue := math.MaxFloat64
	for i, p := range topProcs {
		val := getValue(p, sortMode)
		values[i] = val
		if val > maxValue {
			maxValue = val
		}
		if val < minValue {
			minValue = val
		}
	}

	// Create bubbles with sizes
	bubbles := make([]Bubble, count)
	for i, p := range topProcs {
		bubbles[i] = Bubble{
			Process:  p,
			Value:    values[i],
			Radius:   calculateBubbleSize(values[i], maxValue, minValue, width, height),
			ColorHex: getColor(values[i], maxValue, minValue),
		}
	}

	// Position bubbles using circle packing
	// Start with largest at center
	centerX := float64(width) / 2.0
	centerY := float64(height) / 2.0

	if len(bubbles) > 0 {
		bubbles[0].X = centerX
		bubbles[0].Y = centerY
	}

	// Pack remaining bubbles
	for i := 1; i < len(bubbles); i++ {
		placed := false
		currentRadius := bubbles[i].Radius

		// Try positions in expanding spiral
		for distance := currentRadius + bubbles[0].Radius; distance < float64(width)*2 && !placed; distance += 2 {
			for angle := 0.0; angle < 2*math.Pi; angle += 0.3 {
				x := centerX + distance*math.Cos(angle)
				y := centerY + distance*math.Sin(angle)

				// Check if this position collides with any existing bubble
				collides := false
				for j := 0; j < i; j++ {
					dx := x - bubbles[j].X
					dy := y - bubbles[j].Y
					minDist := currentRadius + bubbles[j].Radius + 1
					if math.Sqrt(dx*dx+dy*dy) < minDist {
						collides = true
						break
					}
				}

				// Check screen bounds
				if x-currentRadius < 0 || x+currentRadius >= float64(width) ||
					y-currentRadius < 3 || y+currentRadius >= float64(height-2) {
					collides = true
				}

				if !collides {
					bubbles[i].X = x
					bubbles[i].Y = y
					placed = true
					break
				}
			}
		}

		// If we couldn't place it, skip this bubble
		if !placed {
			bubbles = bubbles[:i]
			break
		}
	}

	return bubbles
}

// renderBubbles creates a 2D character array and draws all bubbles
func renderBubbles(bubbles []Bubble, width, height int, sortMode string) string {
	// Create empty canvas
	canvas := make([][]rune, height)
	for i := range canvas {
		canvas[i] = make([]rune, width)
		for j := range canvas[i] {
			canvas[i][j] = ' '
		}
	}

	// Store colors for each position
	colors := make([][]string, height)
	for i := range colors {
		colors[i] = make([]string, width)
	}

	// Draw each bubble
	for _, b := range bubbles {
		drawBubble(canvas, colors, b, sortMode)
	}

	// Convert canvas to string with colors
	var result strings.Builder
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			char := canvas[y][x]
			color := colors[y][x]

			if color != "" {
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
				result.WriteString(style.Render(string(char)))
			} else {
				result.WriteRune(char)
			}
		}
		if y < height-1 {
			result.WriteRune('\n')
		}
	}

	return result.String()
}

// drawBubble draws a single bubble on the canvas
func drawBubble(canvas [][]rune, colors [][]string, b Bubble, sortMode string) {
	height := len(canvas)
	width := len(canvas[0])

	// Draw circle outline
	for angle := 0.0; angle < 2*math.Pi; angle += 0.1 {
		x := int(b.X + b.Radius*math.Cos(angle))
		y := int(b.Y + b.Radius*math.Sin(angle))

		if y >= 0 && y < height && x >= 0 && x < width {
			canvas[y][x] = '○'
			colors[y][x] = b.ColorHex
		}
	}

	// Create label
	metricSuffix := ""
	switch sortMode {
	case "cpu":
		metricSuffix = fmt.Sprintf("%.1f%%", b.Value)
	case "memory":
		metricSuffix = fmt.Sprintf("%.1f%%", b.Value)
	case "threads":
		metricSuffix = fmt.Sprintf("%d", int(b.Value))
	}

	label := truncate(b.Process.Name, 12) + " " + metricSuffix

	// Draw label at center
	labelX := int(b.X) - len(label)/2
	labelY := int(b.Y)

	if labelY >= 0 && labelY < height {
		for i, ch := range label {
			x := labelX + i
			if x >= 0 && x < width {
				canvas[labelY][x] = ch
				colors[labelY][x] = b.ColorHex
			}
		}
	}
}
