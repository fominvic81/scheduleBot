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

	file, ferr := os.OpenFile("logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if ferr != nil {
		log.Fatal(ferr)
	}
	log.SetOutput(file)

	database, err := db.Init()

	if err != nil {
		log.Fatal("Failed to init db." + err.Error())
		return
	}

	telegram.Init(os.Getenv("TOKEN"), database)
}
