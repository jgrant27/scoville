package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TargetPrice          float64
	TotalBudget          float64
	PaperMode            bool
	EnableLiquidityPhase bool
	WhaleBuyMin          float64
	WhaleBuyMax          float64
	RetailBuyMin         float64
	RetailBuyMax         float64
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		TargetPrice:          getRequiredFloat("TARGET_PRICE"),
		TotalBudget:          getRequiredFloat("TOTAL_BUDGET"),
		PaperMode:            getRequiredBool("PAPER_MODE"),
		EnableLiquidityPhase: getRequiredBool("ENABLE_LIQUIDITY_PHASE"),
		WhaleBuyMin:          getRequiredFloat("WHALE_BUY_MIN"),
		WhaleBuyMax:          getRequiredFloat("WHALE_BUY_MAX"),
		RetailBuyMin:         getRequiredFloat("RETAIL_BUY_MIN"),
		RetailBuyMax:         getRequiredFloat("RETAIL_BUY_MAX"),
	}
}

func getRequiredEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("FATAL: Missing required environment variable: %s", key)
	}
	return value
}

func getRequiredFloat(key string) float64 {
	valStr := getRequiredEnv(key)
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		log.Fatalf("FATAL: Invalid value for %s: %v", key, err)
	}
	return val
}

func getRequiredBool(key string) bool {
	valStr := getRequiredEnv(key)
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		log.Fatalf("FATAL: Invalid value for %s: %v", key, err)
	}
	return val
}
