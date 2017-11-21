package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/clinstid/envir_collector_go/shared"
)

var dbClient *shared.DBClient

func readingsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request:", r)
	now := time.Now()

	d, err := time.ParseDuration("-1h")
	if err != nil {
		log.Panic("Failed to parse duration:", err)
		return
	}
	startTime := now.Add(d)
	endTime := now

	log.Println("Running query from", startTime, "to", endTime, "...")
	readings, err := dbClient.GetReadings(startTime, endTime)
	if err != nil {
		log.Panic("Failed to get readings:", err)
	}
	readingsJson, err := json.Marshal(&readings)
	if err != nil {
		log.Panic("Failed to marshal JSON from readings", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(readingsJson)
}

func main() {
	// DB configuration parameters
	dbHost := flag.String("db-host", "yoda", "database host name")
	dbUser := flag.String("db-user", "energydash", "Database user")
	dbPassword := flag.String("db-password", "energydash", "Database password")
	dbName := flag.String("db-name", "energydash", "Database name")
	dbClient = shared.NewDBClient(*dbHost, *dbUser, *dbPassword, *dbName)

	log.Println("Starting http server")
	http.HandleFunc("/readings", readingsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
