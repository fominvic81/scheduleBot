package main

import (
	"log"

	"github.com/fominvic81/scheduleBot/db"
	"github.com/fominvic81/scheduleBot/telegram"

	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	database, err := db.Init()

	if err != nil {
		log.Fatal("Failed to init db." + err.Error())
		return
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("Failed to get TOKEN env variable")
	}

	file, ferr := os.Open("logs.txt")
	if ferr != nil {
		log.Fatal(ferr)
	}
	log.SetOutput(file)

	telegram.Init(token, database)
}
