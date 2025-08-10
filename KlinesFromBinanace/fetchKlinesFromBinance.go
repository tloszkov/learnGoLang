package klinesfrombinance

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	loadenv "learnGoLang/LoadEnv"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

type Candle struct {
	Symbol    string
	Timestamp int64 // Unix timestamp in milliseconds
	Datetime  time.Time
	Date      string
	Hour      int
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

func parseCSV(filePath string) ([]Candle, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read() // Read header
	if err == io.EOF {           // Empty file
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Simple check for expected columns (flexible to order)
	colMap := make(map[string]int)
	for i, colName := range header {
		colMap[colName] = i
	}

	expectedCols := []string{"symbol", "timestamp", "datetime", "date", "hour", "open", "high", "low", "close", "volume"}
	for _, ec := range expectedCols {
		if _, ok := colMap[ec]; !ok {
			return nil, fmt.Errorf("missing expected column in CSV: %s", ec)
		}
	}

	var candles []Candle
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		timestamp, _ := strconv.ParseInt(record[colMap["timestamp"]], 10, 64)
		datetime, _ := time.Parse("2006-01-02 15:04:05", record[colMap["datetime"]])
		hour, _ := strconv.Atoi(record[colMap["hour"]])
		open, _ := strconv.ParseFloat(record[colMap["open"]], 64)
		high, _ := strconv.ParseFloat(record[colMap["high"]], 64)
		low, _ := strconv.ParseFloat(record[colMap["low"]], 64)
		close, _ := strconv.ParseFloat(record[colMap["close"]], 64)
		volume, _ := strconv.ParseFloat(record[colMap["volume"]], 64)

		candles = append(candles, Candle{
			Symbol: record[colMap["symbol"]], Timestamp: timestamp, Datetime: datetime,
			Date: record[colMap["date"]], Hour: hour, Open: open, High: high, Low: low,
			Close: close, Volume: volume,
		})
	}
	return candles, nil
}

// fetchKlinesFromBinance fetches kline data from Binance API for a given time range
func fetchKlinesFromBinance(symbol, interval string, startTime, endTime int64) ([]Candle, error) {
	url := fmt.Sprintf("%s?symbol=%s&interval=%s&startTime=%d&endTime=%d&limit=1000",
		loadenv.BINANCE_API_BASE, symbol, interval, startTime, endTime)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Binance API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Binance API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Binance API response: %w", err)
	}

	var rawKlines [][]interface{}
	err = json.Unmarshal(body, &rawKlines)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Binance API response: %w", err)
	}

	var candles []Candle
	for _, kline := range rawKlines {
		if len(kline) < 11 { // Binance kline has at least 11 fields
			continue
		}

		openTime := int64(kline[0].(float64)) // Timestamp in milliseconds
		openPrice, _ := strconv.ParseFloat(kline[1].(string), 64)
		highPrice, _ := strconv.ParseFloat(kline[2].(string), 64)
		lowPrice, _ := strconv.ParseFloat(kline[3].(string), 64)
		closePrice, _ := strconv.ParseFloat(kline[4].(string), 64)
		volume, _ := strconv.ParseFloat(kline[5].(string), 64)

		dt := time.Unix(0, openTime*int64(time.Millisecond))

		candles = append(candles, Candle{
			Symbol:    symbol,
			Timestamp: openTime,
			Datetime:  dt,
			Date:      dt.Format("2006-01-02"),
			Hour:      dt.Hour(),
			Open:      openPrice,
			High:      highPrice,
			Low:       lowPrice,
			Close:     closePrice,
			Volume:    volume,
		})
	}
	return candles, nil
}

// updateHistoricalData fetches and updates CSV with missing data from Binance
func updateHistoricalData(filePath, symbol, interval string) ([]Candle, error) {
	fmt.Printf("Updating data: %s pair, %s interval, file: %s\n", symbol, interval, filePath)

	existingCandles, err := parseCSV(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error reading existing CSV: %w", err)
	}

	var lastTimestamp int64 = 0
	if len(existingCandles) > 0 {
		// Sort candles by timestamp to ensure the last one is truly the latest
		sort.Slice(existingCandles, func(i, j int) bool {
			return existingCandles[i].Timestamp < existingCandles[j].Timestamp
		})
		lastTimestamp = existingCandles[len(existingCandles)-1].Timestamp
		fmt.Printf("Last data in CSV: %s (timestamp: %d)\n", existingCandles[len(existingCandles)-1].Datetime, lastTimestamp)
	} else {
		fmt.Println("NO DATA, STARTING FROM THE BEGINNING (2025 JAN 17)")
		lastTimestamp = time.Date(2025, time.January, 17, 0, 0, 0, 0, time.UTC).UnixMilli()
	}

	// Calculate the duration of one interval in milliseconds for fetching batches
	var intervalDuration time.Duration
	switch interval {
	case "1m":
		intervalDuration = time.Minute
	case "3m":
		intervalDuration = 3 * time.Minute
	case "5m":
		intervalDuration = 5 * time.Minute
	case "15m":
		intervalDuration = 15 * time.Minute
	case "30m":
		intervalDuration = 30 * time.Minute
	case "1h":
		intervalDuration = time.Hour
	case "2h":
		intervalDuration = 2 * time.Hour
	case "4h":
		intervalDuration = 4 * time.Hour
	case "6h":
		intervalDuration = 6 * time.Hour
	case "8h":
		intervalDuration = 8 * time.Hour
	case "12h":
		intervalDuration = 12 * time.Hour
	case "1d":
		intervalDuration = 24 * time.Hour
	case "3d":
		intervalDuration = 3 * 24 * time.Hour
	case "1w":
		intervalDuration = 7 * 24 * time.Hour
	case "1M":
		// This is tricky, 1 month duration varies. For fetching, it's safer to use a fixed max batch size.
		// For '1M', it's best to rethink how 'intervalDurationMillis' is used here, or use smaller intervals.
		// For now, let's assume it's roughly 30 days for batching purposes, but API handles exact month boundaries.
		intervalDuration = 30 * 24 * time.Hour // Approximate for batching
	default:
		fmt.Printf("Warning: Unknown interval '%s', assuming 1-hour duration for batch size.\n", interval)
		intervalDuration = time.Hour
	}
	intervalDurationMillis := int64(intervalDuration.Milliseconds())

	// Fetch data up to current time, truncating to the nearest interval start
	now := time.Now().UTC()
	fetchEndTime := now.Truncate(intervalDuration).UnixMilli()

	if lastTimestamp == 0 { // If no existing data, start from the defined START_DATE_STR
		startTimeObj, err := time.Parse("2025-01-02", "2025-01-02")
		if err != nil {
			return nil, fmt.Errorf("invalid START_DATE_STR format: %w", err)
		}
		// Adjust start time to be aligned with the interval if necessary
		lastTimestamp = startTimeObj.Truncate(intervalDuration).UnixMilli()
		fmt.Printf("Starting date based on START_DATE_STR: %s (timestamp: %d)\n", time.Unix(0, lastTimestamp*int64(time.Millisecond)), lastTimestamp)
	}

	// Adjust lastTimestamp to be the start of the next candle for fetching
	fetchStartTime := lastTimestamp + intervalDurationMillis

	if fetchStartTime >= fetchEndTime {
		fmt.Println("The file is up to date, no need to download new data.")
		return existingCandles, nil
	}

	fmt.Printf("Fetching data from Binance: from %s to %s...\n",
		time.Unix(0, fetchStartTime*int64(time.Millisecond)).Format("2006-01-02 15:04:05"),
		time.Unix(0, fetchEndTime*int64(time.Millisecond)).Format("2006-01-02 15:04:05"))

	var newCandles []Candle
	// Binance limit is 1000 candles per request. Need to loop for larger ranges.
	for currentBatchStartTime := fetchStartTime; currentBatchStartTime < fetchEndTime; {
		// Calculate batchEndTime based on 1000 candles * intervalDuration
		batchEndTime := currentBatchStartTime + (999 * intervalDurationMillis)
		if batchEndTime >= fetchEndTime {
			batchEndTime = fetchEndTime
		}

		fmt.Printf("Fetching: %s -> %s\n",
			time.Unix(0, currentBatchStartTime*int64(time.Millisecond)).Format("2006-01-02 15:04:05"),
			time.Unix(0, batchEndTime*int64(time.Millisecond)).Format("2006-01-02 15:04:05"))

		batchCandles, err := fetchKlinesFromBinance(symbol, interval, currentBatchStartTime, batchEndTime)
		if err != nil {
			return nil, fmt.Errorf("error fetching batch from Binance: %w", err)
		}
		newCandles = append(newCandles, batchCandles...)

		if len(batchCandles) == 0 { // No more data in this range
			break
		}
		// Move to the start of the next candle after the last fetched candle
		currentBatchStartTime = batchCandles[len(batchCandles)-1].Timestamp + intervalDurationMillis
		time.Sleep(100 * time.Millisecond) // Be nice to the API
	}

	if len(newCandles) == 0 {
		fmt.Println("No new data downloaded.")
		return existingCandles, nil
	}

	fmt.Printf("Successfully downloaded %d new candles.\n", len(newCandles))

	// Append new candles to existing ones
	allCandles := append(existingCandles, newCandles...)

	// Remove duplicates (if any) and sort
	uniqueCandlesMap := make(map[int64]Candle)
	for _, c := range allCandles {
		uniqueCandlesMap[c.Timestamp] = c
	}
	allCandles = make([]Candle, 0, len(uniqueCandlesMap))
	for _, c := range uniqueCandlesMap {
		allCandles = append(allCandles, c)
	}
	sort.Slice(allCandles, func(i, j int) bool {
		return allCandles[i].Timestamp < allCandles[j].Timestamp
	})

	// Write all unique and sorted candles back to the CSV
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV for writing: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ',' // Use comma as delimiter

	// Write header
	writer.Write([]string{"symbol", "timestamp", "datetime", "date", "hour", "open", "high", "low", "close", "volume"})

	for _, c := range allCandles {
		record := []string{
			c.Symbol,
			strconv.FormatInt(c.Timestamp, 10),
			c.Datetime.Format("2006-01-02 15:04:05"),
			c.Date,
			strconv.Itoa(c.Hour),
			strconv.FormatFloat(c.Open, 'f', 8, 64),
			strconv.FormatFloat(c.High, 'f', 8, 64),
			strconv.FormatFloat(c.Low, 'f', 8, 64),
			strconv.FormatFloat(c.Close, 'f', 8, 64),
			strconv.FormatFloat(c.Volume, 'f', 8, 64),
		}
		writer.Write(record)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing CSV writer: %w", err)
	}

	fmt.Printf("CSV file successfully updated: %s\n", filePath)
	return allCandles, nil
}

func CreateDataFolder() error {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		return os.Mkdir("data", 0755)
	}
	return nil
}

func FetchData() ([]Candle, error) {
	err := CreateDataFolder()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Data folder created successfully.")
	}
	data, err := updateHistoricalData(loadenv.DATA_FILE_PATH, loadenv.SYMBOL, loadenv.BINANCE_INTERVAL)
	return data, err
}
