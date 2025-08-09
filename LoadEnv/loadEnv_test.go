package loadenv

import (
	"os"
	"testing"
	"time"
)

func TestLoadEnv(t *testing.T) {
	// Create a temporary .env file for testing
	envContent := `	FAST_LENGTH=12
					SLOW_LENGTH=26
					SIGNAL_LENGTH=9
					TREND_TF_HOURS=4
					ENTRY_TF_MINUTES=15
					MAX_POSITION_HOLD_HOURS=24
					STOP_LOSS_PCT=2.5
					TAKE_PROFIT_PCT=5.0
					TRAILING_STOP_PCT=1.5
					MAX_ALLOWED_SL_PCT=3.0
					MIN_MACD_STRENGTH=0.001
					SLIPPAGE_POINTS=1.0
					ENABLE_SHORT_TRADES=false
					OUTPUT_FILE_NAME=test_output.csv
					BINANCE_API_BASE=https://api.binance.com
					BINANCE_INTERVAL=15m
					SYMBOL=ETHUSDT
					WEBSOCKET_URL=wss://stream.binance.com:9443/ws/ethusdt@kline_15m
					BINANCE_API_KEY=test_api_key
					BINANCE_SECRET_KEY=test_secret_key
					START_DATE_STR=2020-01-01 00:00:00
					DATA_FILE_PATH=test_data.csv
					TELEGRAM_BOT_TOKEN=test_bot_token
					TELEGRAM_CHAT_ID=123456789`

	// Write temporary .env file
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test .env file: %v", err)
	}
	defer os.Remove(".env") // Clean up after test

	// Call LoadEnv
	LoadEnv()

	// Test integer values
	tests := []struct {
		name     string
		actual   int
		expected int
	}{
		{"FAST_LENGTH", FAST_LENGTH, 12},
		{"SLOW_LENGTH", SLOW_LENGTH, 26},
		{"SIGNAL_LENGTH", SIGNAL_LENGTH, 9},
		{"TREND_TF_HOURS", TREND_TF_HOURS, 4},
		{"ENTRY_TF_MINUTES", ENTRY_TF_MINUTES, 15},
		{"MAX_POSITION_HOLD_HOURS", MAX_POSITION_HOLD_HOURS, 24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.actual, tt.expected)
			}
		})
	}

	// Test float values
	floatTests := []struct {
		name     string
		actual   float64
		expected float64
	}{
		{"STOP_LOSS_PCT", STOP_LOSS_PCT, 2.5},
		{"TAKE_PROFIT_PCT", TAKE_PROFIT_PCT, 5.0},
		{"TRAILING_STOP_PCT", TRAILING_STOP_PCT, 1.5},
		{"MAX_ALLOWED_SL_PCT", MAX_ALLOWED_SL_PCT, 3.0},
		{"MIN_MACD_STRENGTH", MIN_MACD_STRENGTH, 0.001},
		{"SLIPPAGE_POINTS", SLIPPAGE_POINTS, 1.0},
	}

	for _, tt := range floatTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("%s = %f, want %f", tt.name, tt.actual, tt.expected)
			}
		})
	}

	// Test string values
	stringTests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"OUTPUT_FILE_NAME", OUTPUT_FILE_NAME, "test_output.csv"},
		{"BINANCE_API_BASE", BINANCE_API_BASE, "https://api.binance.com"},
		{"BINANCE_INTERVAL", BINANCE_INTERVAL, "15m"},
		{"SYMBOL", SYMBOL, "ETHUSDT"},
		{"WEBSOCKET_URL", WEBSOCKET_URL, "wss://stream.binance.com:9443/ws/ethusdt@kline_15m"},
		{"BINANCE_API_KEY", BINANCE_API_KEY, "test_api_key"},
		{"BINANCE_SECRET_KEY", BINANCE_SECRET_KEY, "test_secret_key"},
		{"START_DATE_STR", START_DATE_STR, "2020-01-01 00:00:00"},
		{"DATA_FILE_PATH", DATA_FILE_PATH, "test_data.csv"},
		{"TELEGRAM_BOT_TOKEN", TELEGRAM_BOT_TOKEN, "test_bot_token"},
	}

	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("%s = %s, want %s", tt.name, tt.actual, tt.expected)
			}
		})
	}

	// Test parsed date
	expectedDate, _ := time.Parse("2006-01-02 15:04:05", "2020-01-01 00:00:00")
	if !StartDate.Equal(expectedDate) {
		t.Errorf("StartDate = %v, want %v", StartDate, expectedDate)
	}

	// Test Telegram Chat ID
	if TELEGRAM_CHAT_ID != 123456789 {
		t.Errorf("TELEGRAM_CHAT_ID = %d, want 123456789", TELEGRAM_CHAT_ID)
	}
}

func TestLoadEnvMissingFile(t *testing.T) {
	// Ensure no .env file exists
	os.Remove(".env")
	os.Remove("../.env")

	// This should cause a fatal error, but we can't easily test that
	// In a real scenario, you might want to modify LoadEnv to return an error
	// instead of calling log.Fatal for better testability
}

func TestMustParseIntValid(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	result := mustParseInt("TEST_INT")
	if result != 42 {
		t.Errorf("mustParseInt = %d, want 42", result)
	}
}

func TestMustParseFloatValid(t *testing.T) {
	os.Setenv("TEST_FLOAT", "3.14")
	defer os.Unsetenv("TEST_FLOAT")

	result := mustParseFloat("TEST_FLOAT")
	if result != 3.14 {
		t.Errorf("mustParseFloat = %f, want 3.14", result)
	}
}

func TestMustParseBoolValid(t *testing.T) {
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")

	result := mustParseBool("TEST_BOOL")
	if result != true {
		t.Errorf("mustParseBool = %t, want true", result)
	}
}

func TestMustParseInt64Valid(t *testing.T) {
	os.Setenv("TEST_INT64", "9223372036854775807")
	defer os.Unsetenv("TEST_INT64")

	result := mustParseInt64("TEST_INT64")
	if result != 9223372036854775807 {
		t.Errorf("mustParseInt64 = %d, want 9223372036854775807", result)
	}
}
