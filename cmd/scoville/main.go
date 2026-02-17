package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chili-network/scoville/internal/config"
	"github.com/chili-network/scoville/internal/engine"
	"github.com/chili-network/scoville/internal/state"
	"github.com/chili-network/scoville/internal/ui"
)

func generateResumeLogs(s state.Progress) []string {
	if s.CurrentPhase == 1 && s.WhaleIndex == 0 && s.RetailIndex == 0 && s.TotalSpent == 0 {
		return []string{"Initializing Scoville TUI..."}
	}

	logs := []string{"Resuming previous session from scoville_progress.json..."}

	if s.CurrentPhase > 3 {
		logs = append(logs, "✅ Mission was already complete.")
		logs = append(logs, fmt.Sprintf("Total spent: $%.2f", s.TotalSpent))
		return logs
	}

	logs = append(logs, fmt.Sprintf("Current phase: %d", s.CurrentPhase))
	logs = append(logs, fmt.Sprintf("Whale buys completed: %d", s.WhaleIndex))
	logs = append(logs, fmt.Sprintf("Retail buys completed: %d", s.RetailIndex))
	logs = append(logs, fmt.Sprintf("Total spent so far: $%.2f", s.TotalSpent))

	if s.IsPaused {
		logs = append(logs, "Bot is currently PAUSED. Press 'P' to resume.")
	}

	return logs
}

func main() {
	cfg := config.Load()

	if !cfg.PaperMode {
		fmt.Println("====================================================================")
		fmt.Println("  WARNING: LIVE TRADING MODE IS ENABLED")
		fmt.Println("====================================================================")
		fmt.Println("This will execute REAL trades on the blockchain with REAL funds.")
		fmt.Println("Ensure your configuration in the .env file is correct before proceeding.")
		fmt.Println("")
		fmt.Print("Type 'YES' to confirm and start live trading: ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input != "YES" {
			fmt.Println("Confirmation not received. Exiting.")
			os.Exit(0)
		}
		fmt.Println("Confirmation received. Starting Scoville...")
	}

	s := state.Load()
	initialLogs := generateResumeLogs(s)

	// Initialize Bot
	bot := engine.New(cfg, s)

	// Initialize TUI
	p := tea.NewProgram(ui.InitialModel(bot, s, initialLogs), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
