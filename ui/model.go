package ui

import (
	"fmt"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/req/procmap/proc"
)

type tickMsg time.Time

type model struct {
	processes  []proc.ProcessInfo
	cursor     int
	height     int
	width      int
	err        error
	sortMode   string // "cpu", "memory", "threads"
	maxBubbles int    // configurable number of bubbles to show
}

func NewModel() model {
	return model{
		sortMode:   "cpu",
		maxBubbles: 20,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), fetchProcesses)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func fetchProcesses() tea.Msg {
	procs, err := proc.GetProcesses()
	if err != nil {
		return errMsg{err}
	}

	return processesMsg(procs)
}

type processesMsg []proc.ProcessInfo
type errMsg struct{ err error }

func (m *model) sortProcesses() {
	switch m.sortMode {
	case "cpu":
		sort.Slice(m.processes, func(i, j int) bool {
			return m.processes[i].CPUPercent > m.processes[j].CPUPercent
		})
	case "memory":
		sort.Slice(m.processes, func(i, j int) bool {
			return m.processes[i].MemPercent > m.processes[j].MemPercent
		})
	case "threads":
		sort.Slice(m.processes, func(i, j int) bool {
			return m.processes[i].NumThreads > m.processes[j].NumThreads
		})
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c":
			m.sortMode = "cpu"
			m.sortProcesses()
		case "m":
			m.sortMode = "memory"
			m.sortProcesses()
		case "t":
			m.sortMode = "threads"
			m.sortProcesses()
		case "+", "=":
			m.maxBubbles += 10
		case "-", "_":
			if m.maxBubbles > 10 {
				m.maxBubbles -= 10
			}
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tickMsg:
		return m, tea.Batch(tick(), fetchProcesses)

	case processesMsg:
		m.processes = msg
		// Sort based on current mode
		m.sortProcesses()
		// Keep cursor in bounds
		if m.cursor >= len(m.processes) {
			m.cursor = len(m.processes) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}

	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err)
	}

	var s string

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("235")).
		Width(m.width)

	modeLabel := "CPU%"
	switch m.sortMode {
	case "memory":
		modeLabel = "Memory%"
	case "threads":
		modeLabel = "Threads"
	}

	header := fmt.Sprintf(" Process Bubbles - Sorted by %s | Showing: %d", modeLabel, m.maxBubbles)
	s += headerStyle.Render(header) + "\n"

	// Bubble visualization
	bubbles := packBubbles(m.processes, m.sortMode, m.maxBubbles, m.width, m.height-3)
	bubbleView := renderBubbles(bubbles, m.width, m.height-3, m.sortMode)
	s += bubbleView + "\n"

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))
	footer := fmt.Sprintf("Total: %d | [c]PU [m]emory [t]hreads | +/- adjust count | [q]uit", len(m.processes))
	s += footerStyle.Render(footer)

	return s
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
