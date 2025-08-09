package main

import (
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
	log.Println("Started date:", time.Now().String())
}
