package engine

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/chili-network/scoville/internal/config"
	"github.com/chili-network/scoville/internal/state"
)

// UIUpdate sends data back to the TUI
type UIUpdate struct {
	LogLine      string
	Progress     state.Progress
	PhaseTotal   int // How many steps in current phase?
	PhaseCurrent int // Current step in phase
}

type Bot struct {
	cfg      *config.Config
	state    state.Progress
	UpdateCh chan UIUpdate
	PauseCh  chan bool // Receives pause toggles
	StopCh   chan bool // Receives stop signal
}

func New(cfg *config.Config, s state.Progress) *Bot {
	return &Bot{
		cfg:      cfg,
		state:    s,
		UpdateCh: make(chan UIUpdate, 10),
		PauseCh:  make(chan bool),
		StopCh:   make(chan bool),
	}
}

func (b *Bot) Run() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Phase Constants
	whaleCount := 3
	retailCount := 15
	if !b.cfg.EnableLiquidityPhase {
		retailCount += 5 // Reallocate LP budget to more retail buys
	}

	// --- PHASE 1: WHALES ---
	if b.state.CurrentPhase == 1 {
		b.log("🐋 Entering Phase 1: Whale Anchors")
		for i := b.state.WhaleIndex; i < whaleCount; i++ {
			b.checkControlSignals()

			whaleRange := b.cfg.WhaleBuyMax - b.cfg.WhaleBuyMin
			amount := b.cfg.WhaleBuyMin + r.Float64()*whaleRange
			b.simulateTrade(fmt.Sprintf("Whale_%d", i+1), amount)

			// Update State
			b.state.WhaleIndex = i + 1
			b.state.TotalSpent += amount
			state.Save(b.state)

			// Emit Update
			b.UpdateCh <- UIUpdate{
				LogLine:      fmt.Sprintf("Whale %d bought $%.2f", i+1, amount),
				Progress:     b.state,
				PhaseTotal:   whaleCount,
				PhaseCurrent: i + 1,
			}
			time.Sleep(1500 * time.Millisecond) // Sim delay
		}
		b.state.CurrentPhase = 2
		state.Save(b.state)
	}

	// --- PHASE 2: RETAIL ---
	if b.state.CurrentPhase == 2 {
		b.log("🐟 Entering Phase 2: Retail Gap Fill")
		for i := b.state.RetailIndex; i < retailCount; i++ {
			b.checkControlSignals()

			retailRange := b.cfg.RetailBuyMax - b.cfg.RetailBuyMin
			amount := b.cfg.RetailBuyMin + r.Float64()*retailRange
			b.simulateTrade(fmt.Sprintf("Retail_%d", i+1), amount)

			b.state.RetailIndex = i + 1
			b.state.TotalSpent += amount
			state.Save(b.state)

			b.UpdateCh <- UIUpdate{
				LogLine:      fmt.Sprintf("Retail %d bought $%.2f", i+1, amount),
				Progress:     b.state,
				PhaseTotal:   retailCount,
				PhaseCurrent: i + 1,
			}
			time.Sleep(500 * time.Millisecond)
		}

		if b.cfg.EnableLiquidityPhase {
			b.state.CurrentPhase = 3
			state.Save(b.state)
		} else {
			b.state.CurrentPhase = 4 // Done
			state.Save(b.state)
			b.log("✅ MISSION COMPLETE.")
		}
	}

	// --- PHASE 3: LP ---
	if b.state.CurrentPhase == 3 {
		b.checkControlSignals()
		b.log("🔒 Entering Phase 3: Liquidity Injection")

		b.UpdateCh <- UIUpdate{
			LogLine:      "Adding Liquidity to Pool...",
			Progress:     b.state,
			PhaseTotal:   1,
			PhaseCurrent: 1,
		}
		time.Sleep(2 * time.Second)

		b.log("✅ LIQUIDITY ADDED. MISSION COMPLETE.")
		b.state.CurrentPhase = 4 // Done
		state.Save(b.state)
	}
}

func (b *Bot) checkControlSignals() {
	select {
	case paused := <-b.PauseCh:
		b.state.IsPaused = paused
		state.Save(b.state)
		if paused {
			b.log("⏸️  PROCESS PAUSED. Waiting for resume...")
			// Block until resumed
			for p := range b.PauseCh {
				if !p {
					b.log("▶️  RESUMED.")
					b.state.IsPaused = false
					break
				}
			}
		}
	case <-b.StopCh:
		b.log("🛑 STOP SIGNAL RECEIVED.")
		panic("GracefulExit") // Handled in main
	default:
		// Continue
	}
}

func (b *Bot) simulateTrade(wallet string, amount float64) {
	// Logic stub
}

func (b *Bot) log(msg string) {
	b.UpdateCh <- UIUpdate{LogLine: msg, Progress: b.state}
}
