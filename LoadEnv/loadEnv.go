package loadenv

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Global variables for strategy parameters
var (
	FAST_LENGTH             int
	SLOW_LENGTH             int
	SIGNAL_LENGTH           int
	TREND_TF_HOURS          int
	ENTRY_TF_MINUTES        int
	STOP_LOSS_PCT           float64
	TAKE_PROFIT_PCT         float64
	TRAILING_STOP_PCT       float64
	MAX_ALLOWED_SL_PCT      float64
	MIN_MACD_STRENGTH       float64
	REQUIRE_CONFIRMATION    bool
	COMMISSION_PERCENT      float64
	SLIPPAGE_POINTS         float64
	MAX_POSITION_HOLD_HOURS int
	OUTPUT_FILE_NAME        string
	ENABLE_SHORT_TRADES     bool
)

// Global variables for Binance API and WebSocket constants
var (
	BINANCE_API_BASE   string
	BINANCE_INTERVAL   string
	SYMBOL             string
	WEBSOCKET_URL      string
	BINANCE_API_KEY    string
	BINANCE_SECRET_KEY string
)

// Global variables for data paths and start date
var (
	START_DATE_STR string
	DATA_FILE_PATH string
	StartDate      time.Time
)

// Optional: Telegram notification variables
var (
	TELEGRAM_BOT_TOKEN string
	TELEGRAM_CHAT_ID   int64
)

// LoadEnv loads environment variables from a .env file and assigns them to global variables.
func LoadEnv() {
	// Attempt to load .env from the current directory, or one level up (where main.go might be)
	err := godotenv.Load() // Loads from ./.env by default
	if err != nil {
		log.Printf("Could not load .env file from current directory, trying parent: %v", err)
		err = godotenv.Load("../.env") // Try loading from parent directory
		if err != nil {
			log.Fatalf("Error loading .env file from current or parent directory: %v", err)
		}
	}

	// Strategy Parameters
	FAST_LENGTH = mustParseInt("FAST_LENGTH")
	SLOW_LENGTH = mustParseInt("SLOW_LENGTH")
	SIGNAL_LENGTH = mustParseInt("SIGNAL_LENGTH")
	TREND_TF_HOURS = mustParseInt("TREND_TF_HOURS")
	ENTRY_TF_MINUTES = mustParseInt("ENTRY_TF_MINUTES")
	MAX_POSITION_HOLD_HOURS = mustParseInt("MAX_POSITION_HOLD_HOURS")

	STOP_LOSS_PCT = mustParseFloat("STOP_LOSS_PCT")
	TAKE_PROFIT_PCT = mustParseFloat("TAKE_PROFIT_PCT")
	TRAILING_STOP_PCT = mustParseFloat("TRAILING_STOP_PCT")
	MAX_ALLOWED_SL_PCT = mustParseFloat("MAX_ALLOWED_SL_PCT")
	MIN_MACD_STRENGTH = mustParseFloat("MIN_MACD_STRENGTH")
	SLIPPAGE_POINTS = mustParseFloat("SLIPPAGE_POINTS")

	OUTPUT_FILE_NAME = os.Getenv("OUTPUT_FILE_NAME")
	if OUTPUT_FILE_NAME == "" {
		log.Fatal("OUTPUT_FILE_NAME not set in .env")
	}

	// Binance API and WebSocket Constants
	BINANCE_API_BASE = os.Getenv("BINANCE_API_BASE")
	if BINANCE_API_BASE == "" {
		log.Fatal("BINANCE_API_BASE not set in .env")
	}
	BINANCE_INTERVAL = os.Getenv("BINANCE_INTERVAL")
	if BINANCE_INTERVAL == "" {
		log.Fatal("BINANCE_INTERVAL not set in .env")
	}
	SYMBOL = os.Getenv("SYMBOL")
	if SYMBOL == "" {
		log.Fatal("SYMBOL not set in .env")
	}
	WEBSOCKET_URL = os.Getenv("WEBSOCKET_URL")
	if WEBSOCKET_URL == "" {
		log.Fatal("WEBSOCKET_URL not set in .env")
	}

	BINANCE_API_KEY = os.Getenv("BINANCE_API_KEY")
	BINANCE_SECRET_KEY = os.Getenv("BINANCE_SECRET_KEY")

	// Global variables for data paths and start date
	START_DATE_STR = os.Getenv("START_DATE_STR")
	if START_DATE_STR == "" {
		START_DATE_STR = "2020-01-02 15:04:05" // Default if not set
	}
	DATA_FILE_PATH = os.Getenv("DATA_FILE_PATH")
	if DATA_FILE_PATH == "" {
		DATA_FILE_PATH = "data/ETHUSDC_15m.csv" // Default if not set
	}
	var parseErr error
	StartDate, parseErr = time.Parse("2006-01-02 15:04:05", START_DATE_STR)

	if parseErr != nil {
		log.Fatalf("Error parsing START_DATE_STR from .env: %v", parseErr)
	}

	// Telegram Notification (Optional)
	telegramChatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	if telegramChatIDStr != "" {
		var err error
		TELEGRAM_CHAT_ID, err = strconv.ParseInt(telegramChatIDStr, 10, 64)
		if err != nil {
			log.Fatalf("Invalid value for TELEGRAM_CHAT_ID in .env: %v", err)
		}
	}

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramBotToken != "" {
		TELEGRAM_BOT_TOKEN = telegramBotToken
	}
}

// Helper functions to parse environment variables or fatal error
func mustParseInt(key string) int {
	s := os.Getenv(key)
	if s == "" {
		log.Fatalf("%s not set in .env", key)
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Invalid value for %s in .env: %v", key, err)
	}
	return v
}

func mustParseInt64(key string) int64 {
	s := os.Getenv(key)
	if s == "" {
		log.Fatalf("%s not set in .env", key)
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalf("Invalid value for %s in .env: %v", key, err)
	}
	return v
}

func mustParseFloat(key string) float64 {
	s := os.Getenv(key)
	if s == "" {
		log.Fatalf("%s not set in .env", key)
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("Invalid value for %s in .env: %v", key, err)
	}
	return v
}

func mustParseBool(key string) bool {
	s := os.Getenv(key)
	if s == "" {
		log.Fatalf("%s not set in .env", key)
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatalf("Invalid value for %s in .env: %v", key, err)
	}
	return v
}
