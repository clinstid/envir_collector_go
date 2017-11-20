package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/clinstid/envir_collector_go/shared"
)

func getenvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); len(value) > 0 {
		if valueAsInt, err := strconv.Atoi(value); err == nil {
			return valueAsInt
		}
		log.Println("WARNING: Failed to convert environment variable", key, "with value", value, "to an int.")
		return fallback
	}
	return fallback
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); len(value) > 0 {
		return value
	}
	return fallback
}

func main() {
	// USB serial device parameters
	bitRate := flag.Int("bit-rate", 57600, "serial bit rate in bps")
	dataBits := flag.Int("data-bits", 8, "serial data bits")
	stopBits := flag.Int("stop-bits", 1, "serial stop bits")
	usbDevice := flag.String("usb-device", "/dev/ttyUSB0", "USB device path")

	// DB configuration parameters
	dbHost := flag.String("db-host", "yoda", "database host name")
	dbUser := flag.String("db-user", "energydash", "Database user")
	dbPassword := flag.String("db-password", "energydash", "Database password")
	dbName := flag.String("db-name", "energydash", "Database name")

	flag.Parse()

	envirClient, err := shared.NewEnvirClient(*bitRate, *dataBits, *stopBits, *usbDevice)
	if err != nil {
		log.Panic("Failed to create an EnvirClient:", err)
		return
	}

	// Create a database client for the writer to use
	dbClient := shared.NewDBClient(*dbHost, *dbUser, *dbPassword, *dbName)

	var c = make(chan shared.XMLMessage, 1000)

	// Start the reader process
	go readValues(envirClient, c)

	// Start the writer process
	go writeValues(c, dbClient)

	var input string
	fmt.Printf("Press enter to stop execution...")
	fmt.Scanln(&input)
}
