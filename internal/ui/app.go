package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/chili-network/scoville/internal/engine"
	"github.com/chili-network/scoville/internal/state"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	titleStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Bold(true).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(subtle).
			Padding(0, 1)

	subtleStyle = lipgloss.NewStyle().Foreground(subtle)
)

type model struct {
	bot             *engine.Bot
	progress        progress.Model
	viewport        viewport.Model
	logs            []string
	quitting        bool
	confirming      bool
	isPaused        bool
	missionComplete bool
	width           int
	height          int

	// Data
	currentPhase int
	phaseCount   int
	phaseTotal   int
}

func InitialModel(bot *engine.Bot, s state.Progress, initialLogs []string) model {
	vp := viewport.New(0, 0)
	vp.YPosition = 0

	return model{
		bot:             bot,
		progress:        progress.New(progress.WithDefaultGradient()),
		viewport:        vp,
		logs:            initialLogs,
		isPaused:        s.IsPaused,
		missionComplete: s.CurrentPhase > 3,
		currentPhase:    s.CurrentPhase,
	}
}

func (m model) Init() tea.Cmd {
	// Start the bot in a goroutine
	go func() {
		defer func() { recover() }() // Catch "GracefulExit" panic
		m.bot.Run()
	}()
	return waitForUpdate(m.bot.UpdateCh)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.confirming {
			if msg.String() == "y" || msg.String() == "Y" {
				return m, tea.Quit
			}
			if msg.String() == "n" || msg.String() == "N" || msg.Type == tea.KeyEsc {
				m.confirming = false
				m.logs = append(m.logs, "Quit cancelled.")
				m.viewport.SetContent(strings.Join(m.logs, "\n"))
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "Q":
			m.confirming = true
			return m, nil
		case "p", "P":
			// Toggle Pause
			m.isPaused = !m.isPaused
			m.bot.PauseCh <- m.isPaused
			status := "PAUSED"
			if !m.isPaused {
				status = "RESUMED"
			}
			m.logs = append(m.logs, fmt.Sprintf(">>> PROCESS %s <<<", status))
			m.viewport.SetContent(strings.Join(m.logs, "\n"))
			return m, nil
		}

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 4

		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)

		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height / 2 // Bottom half
		return m, cmd

	case engine.UIUpdate:
		if msg.LogLine != "" {
			wasAtBottom := m.viewport.AtBottom()
			m.logs = append(m.logs, msg.LogLine)
			// Keep log size manageable
			if len(m.logs) > 100 {
				m.logs = m.logs[len(m.logs)-100:]
			}
			m.viewport.SetContent(strings.Join(m.logs, "\n"))
			if wasAtBottom {
				m.viewport.GotoBottom()
			}
		}

		m.isPaused = msg.Progress.IsPaused
		m.currentPhase = msg.Progress.CurrentPhase
		m.phaseCount = msg.PhaseCurrent
		m.phaseTotal = msg.PhaseTotal

		isPhase3Complete := m.currentPhase == 3 && m.phaseTotal > 0 && m.phaseCount >= m.phaseTotal
		isMissionFinished := m.currentPhase > 3 // Phase 4 means done

		if !m.missionComplete && (isPhase3Complete || isMissionFinished) {
			m.missionComplete = true
			m.logs = append(m.logs, "✅ Mission Complete. Press Q to quit.")
			m.viewport.SetContent(strings.Join(m.logs, "\n"))
			m.viewport.GotoBottom()
		}

		var cmd tea.Cmd
		// Update Progress Bar
		if msg.PhaseTotal > 0 {
			pct := float64(msg.PhaseCurrent) / float64(msg.PhaseTotal)
			cmd = m.progress.SetPercent(pct)
		}

		if m.missionComplete {
			return m, cmd
		}

		return m, tea.Batch(cmd, waitForUpdate(m.bot.UpdateCh))

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.confirming {
		return fmt.Sprintf("\n\n  %s\n\n", titleStyle.Render("Are you sure you want to quit? (Y/N)"))
	}

	// 1. Top Half: Status Dashboard
	phaseName := "Initializing"
	if m.missionComplete {
		phaseName = "🏁 MISSION COMPLETE"
	} else if m.currentPhase == 1 {
		phaseName = "🐋 PHASE 1: WHALES"
	} else if m.currentPhase == 2 {
		phaseName = "🐟 PHASE 2: RETAIL"
	} else if m.currentPhase == 3 {
		phaseName = "🔒 PHASE 3: LIQUIDITY"
	}

	statusView := boxStyle.Render(fmt.Sprintf(
		"%s\n\n%s\n%s",
		titleStyle.Render(phaseName),
		fmt.Sprintf("Sub Count: %d / %d", m.phaseCount, m.phaseTotal),
		m.progress.View(),
	))

	// 2. Bottom Half: Logs
	logView := boxStyle.Render(m.viewport.View())

	// 3. Footer
	help := subtleStyle.Render("Q: Quit | P: Pause/Resume | ↑/↓: Scroll")
	if m.isPaused {
		help = subtleStyle.Render("Q: Quit | P: Resume | PAUSED | ↑/↓: Scroll")
	}

	return lipgloss.JoinVertical(lipgloss.Left, statusView, logView, help)
}

func waitForUpdate(sub chan engine.UIUpdate) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}
