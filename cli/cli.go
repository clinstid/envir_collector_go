package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/clinstid/envir_collector_go/shared"
)

func main() {
	// DB query parameters
	startStr := flag.String("start", "", "Start time in RFC3339 format, default is 1 hour ago")
	endStr := flag.String("end", "", "End time in RFC3339 format, default is now")

	// DB configuration parameters
	dbHost := flag.String("db-host", "yoda", "database host name")
	dbUser := flag.String("db-user", "energydash", "Database user")
	dbPassword := flag.String("db-password", "energydash", "Database password")
	dbName := flag.String("db-name", "energydash", "Database name")

	flag.Parse()

	var err error
	var startTime, endTime time.Time

	now := time.Now()

	if *startStr == "" {
		d, err := time.ParseDuration("-1h")
		if err != nil {
			log.Panic("Failed to parse duration:", err)
			return
		}
		startTime = now.Add(d)
	} else {
		startTime, err = time.Parse(time.RFC3339, *startStr)
		if err != nil {
			log.Panic("Failed to parse start time:", err)
			return
		}
	}

	if *endStr == "" {
		endTime = now
	} else {
		endTime, err = time.Parse(time.RFC3339, *endStr)
		if err != nil {
			log.Panic("Failed to parse end time:", err)
		}
		return
	}

	log.Println("Running query from", startTime, "to", endTime, "...")
	dbClient := shared.NewDBClient(*dbHost, *dbUser, *dbPassword, *dbName)
	readings, err := dbClient.GetReadings(startTime, endTime)
	if err != nil {
		log.Panic("Failed to get readings:", err)
	}
	readingsJson, err := json.Marshal(&readings)
	if err != nil {
		log.Panic("Failed to marshal JSON from readings", err)
	}
	log.Println(string(readingsJson))
}
