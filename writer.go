package main

import (
	"database/sql"
	"fmt"
	"log"
)

func writeValues(c chan XMLMessage, dbHost, dbUser, dbPassword, dbName string) {
	log.Println("[w]: Waiting for messages...")
	dbOpenString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", dbHost, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", dbOpenString)
	if err != nil {
		log.Panic("[w]: Failed to connect to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic("Failed to ping db:", err)
	}

	for {
		message := <-c
		log.Printf("[w]: XMLMessage object: %+v\n", message)
		dbMessage := buildDBMessage(message)
		insertCmd := genInsert(dbMessage)
		log.Printf("[w]: Insert command: %+v\n", insertCmd)
		res, err := db.Exec(insertCmd)
		if err != nil {
			log.Panic("Failed INSERT:", err)
		}
		log.Printf("[w]: Insert result: %+v\n", res)
	}
}
