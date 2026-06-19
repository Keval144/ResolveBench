package tui

import (
	"fmt"
	"sync"
	"time"
	"dnsbench/internal/dns"
	"dnsbench/internal/models"
	"dnsbench/internal/network"

	tea "github.com/charmbracelet/bubbletea"
)

type benchDoneMsg struct {
	results     []models.ProviderResult
	networkRezs []models.NetworkResult
}
type progressUpdateMsg struct {
	current int
	total   int
	label   string
}

type Model struct {
	Program     *tea.Program
	screen      models.Screen
	progress    int
	progressMax int
	progressLbl string
	results     []models.ProviderResult
	networkRezs []models.NetworkResult
	width       int
	height      int
	err         error
	mu          sync.Mutex
}

func NewModel() Model {
	return Model{
		screen: models.WelcomeScreen,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.screen == models.WelcomeScreen {
				m.screen = models.RunningScreen
				m.progress = 0
				m.totalTasks()
				go m.runAll()
				return m, nil
			}
			if m.screen == models.ResultsScreen {
				m.screen = models.NetworkScreen
				return m, nil
			}
		if m.screen == models.NetworkScreen {
			m.screen = models.ResultsScreen
			return m, nil
		}
		case tea.KeySpace:
			if m.screen == models.ResultsScreen || m.screen == models.NetworkScreen {
				m.screen = models.WelcomeScreen
			}
		case tea.KeyEscape, tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyRunes:
			if string(msg.Runes) == "q" {
				return m, tea.Quit
			}
		}

	case progressUpdateMsg:
		m.progress = msg.current
		m.progressMax = msg.total
		m.progressLbl = msg.label
		return m, nil

	case benchDoneMsg:
		m.results = msg.results
		m.networkRezs = msg.networkRezs
		m.screen = models.ResultsScreen
		m.progress = m.progressMax
		return m, nil
	}

	return m, nil
}

func (m *Model) totalTasks() {
	m.progressMax = len(dns.Providers)*len(dns.BenchDomains) + len(dns.NetworkTargets)
}

func (m *Model) runAll() {
	totalMax := m.progressMax
	netOffset := len(dns.Providers) * len(dns.BenchDomains)

	prog := func(current int, label string) {
		if m.Program != nil {
			m.Program.Send(progressUpdateMsg{current: current, total: totalMax, label: label})
		}
	}

	results := dns.RunBenchmark(func(current, total int, label string) {
		prog(current, "DNS: "+label)
	})

	prog(netOffset, "Network tests")

	netResults := network.RunLatencyTests(func(current, total int, label string) {
		prog(netOffset+current, "Net: "+label)
	})

	if m.Program != nil {
		m.Program.Send(benchDoneMsg{results: results, networkRezs: netResults})
	}
}

func (m *Model) View() string {
	switch m.screen {
	case models.WelcomeScreen:
		return welcomeView()
	case models.RunningScreen:
		return m.runningView()
	case models.ResultsScreen:
		return m.resultsView()
	case models.NetworkScreen:
		return m.networkView()
	default:
		return "unknown screen"
	}
}

func formatDur(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(d.Microseconds()))
	}
	return fmt.Sprintf("%.2fms", float64(d.Milliseconds()))
}
