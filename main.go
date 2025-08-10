package main

import (
	klinesfrombinance "learnGoLang/KlinesFromBinanace"
	loadenv "learnGoLang/LoadEnv"
	"log"
	"os"
	"time"
)

func main() {
	loadenv.LoadEnv()
	log.Println("Environment variables loaded successfully.")
	log.Println("Symbol:", loadenv.SYMBOL)
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Start date:", time.Now().Format("2006-01-02 15:04:05"))
	// message := fmt.Sprintf("Ro Bot service started! Symbol:%s | Time:%s", loadenv.SYMBOL, time.Now().Format("2006-01-02 15:04:05"))
	// sendnotification.SendTelegramNotification(message)
	_, err := klinesfrombinance.FetchData()
	if err != nil {
		log.Printf("Error updating historical data: %v\n", err)
		return
	} else {
		log.Println("Historical data updated successfully.")
	}
}
