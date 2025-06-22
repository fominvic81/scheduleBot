package main

import (
	"log"

	"github.com/fominvic81/scheduleBot/db"
	"github.com/fominvic81/scheduleBot/telegram"

	"os"

	"github.com/joho/godotenv"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recover: ", r)
		}
	}()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	database, err := db.Init()

	if err != nil {
		log.Fatal("Failed to init db." + err.Error())
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("Failed to get TOKEN env variable")
	}

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	telegram.Init(token, database)
}
