package config

import (
"os"
"strconv"
"github.com/joho/godotenv"
)

type Config struct {
TargetPrice   float64
TotalBudget   float64
PaperMode     bool
}

func Load() *Config {
_ = godotenv.Load()

price, _ := strconv.ParseFloat(getEnv("TARGET_PRICE", "0.00075"), 64)
budget, _ := strconv.ParseFloat(getEnv("TOTAL_BUDGET", "3000.0"), 64)
paper, _ := strconv.ParseBool(getEnv("PAPER_MODE", "true"))

return &Config{price, budget, paper}
}

func getEnv(key, fallback string) string {
if value, ok := os.LookupEnv(key); ok { return value }
return fallback
}
