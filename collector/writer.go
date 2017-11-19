package main

import (
	"log"

	"github.com/clinstid/envir_collector_go/shared"
)

func writeValues(c chan shared.XMLMessage, dbClient *shared.DBClient) {
	log.Println("[w]: Waiting for messages...")

	for {
		message := <-c
		err := dbClient.WriteMessage(message)
		if err != nil {
			log.Panic("Failed to write message", message, "to database:", err)
		}
	}
}
