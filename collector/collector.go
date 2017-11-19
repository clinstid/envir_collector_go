package main

import (
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
	// TODO Load configuration using flags instead of env vars
	// Load USB configuration from env
	bitRate := getenvAsInt("ENVIR_SERIAL_BIT_RATE", 57600)
	dataBits := getenvAsInt("ENVIR_SERIAL_DATA_BITS", 8)
	stopBits := getenvAsInt("ENVIR_SERIAL_STOP_BITS", 1)
	usbDevice := getenv("ENVIR_SERIAL_USB_DEVICE", "/dev/ttyUSB0")

	envirClient, err := shared.NewEnvirClient(bitRate, dataBits, stopBits, usbDevice)
	if err != nil {
		log.Panic("Failed to create an EnvirClient:", err)
		return
	}

	// Load DB configuration from env
	dbHost := getenv("ENVIR_DB_HOST", "yoda") // TODO: Switch to localhost
	dbUser := getenv("ENVIR_DB_USER", "energydash")
	dbPassword := getenv("ENVIR_DB_PASSWORD", "energydash")
	dbName := getenv("ENVIR_DB_NAME", "energydash")

	// Create a database client for the writer to use
	dbClient := shared.NewDBClient(dbHost, dbUser, dbPassword, dbName)

	var c = make(chan shared.XMLMessage, 1000)

	// Start the reader process
	go readValues(envirClient, c)

	// Start the writer process
	go writeValues(c, dbClient)

	var input string
	fmt.Printf("Press enter to stop execution...")
	fmt.Scanln(&input)
}
