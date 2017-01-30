package main

import (
	"encoding/xml"
	"fmt"
	"github.com/mikepb/go-serial"
	"log"
	"os"
	"strconv"
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

// Channel XML element
type Channel struct {
	Watts int `xml:"watts"`
}

// Message XML element
type Message struct {
	Src    string  `xml:"src"`
	Dsb    string  `xml:"dsb"`
	Time   string  `xml:"time"`
	TmprF  float32 `xml:"tmprF"`
	Sensor int     `xml:"sensor"`
	ID     string  `xml:"id"`
	Type   int     `xml:"type"`
	Ch1    Channel `xml:"ch1"`
	Ch2    Channel `xml:"ch2"`
	Ch3    Channel `xml:"ch3"`
}

func readValues(bitRate, dataBits, stopBits int, usbDevice string, c chan Message) {
	// Reading Example:
	/*
		exampleXML := `
		<msg>
		    <src>CC128-v0.15</src>
		    <dsb>01331</dsb>
		    <time>12:43:53</time>
		    <tmprF>73.5</tmprF>
		    <sensor>0</sensor>
		    <id>00077</id>
		    <type>1</type>
		    <ch1>
		        <watts>00072</watts>
		    </ch1>
		    <ch2>
		        <watts>00637</watts>
		    </ch2>
		    <ch3>
		        <watts>01189</watts>
		    </ch3>
		</msg>
		`
		var exampleMessage Message
		err := xml.Unmarshal([]byte(exampleXml), &exampleMessage)
		if err != nil {
			log.Panic("Failed on example xml:", err)
		}
		log.Println("Example XML result:", exampleMessage)
	*/
	options := serial.Options{BitRate: bitRate, DataBits: dataBits, StopBits: stopBits}
	log.Println("[r]: Connecting to device")
	p, err := options.Open(usbDevice)
	if err != nil {
		log.Panic(err)
	}

	defer p.Close()

	log.Println("[r]: Waiting for data...")

	msg := make([]byte, 0)
	buf := make([]byte, 1)
	for {
		if _, err := p.Read(buf); err != nil {
			log.Panic("Read failed", err)
		} else {
			msg = append(msg[:], buf[0])
			if buf[0] == '\n' {
				msgStr := string(msg[:])
				log.Println("[r]: Message bytes:", msg)
				log.Println("[r]: Message string:", msgStr)

				var message Message
				err = xml.Unmarshal(msg, &message)
				if err != nil {
					log.Panic("[r]: xml.Unmarshal failed:", err)
				} else {
					c <- message
				}

				msg = make([]byte, 0)
			}
		}
	}
}

func writeValues(c chan Message) {
	log.Println("[w]: Waiting for messages...")
	for {
		message := <-c
		log.Printf("[w]: Message object: %+v\n", message)
	}
}

func main() {
	bitRate := getenvAsInt("ENVIR_SERIAL_BIT_RATE", 57600)
	dataBits := getenvAsInt("ENVIR_SERIAL_DATA_BITS", 8)
	stopBits := getenvAsInt("ENVIR_SERIAL_STOP_BITS", 1)
	usbDevice := getenv("ENVIR_SERIAL_USB_DEVICE", "/dev/ttyUSB0")

	var c chan Message = make(chan Message, 1000)

	go readValues(bitRate, dataBits, stopBits, usbDevice, c)
	go writeValues(c)

	var input string
	fmt.Printf("Press enter to stop execution...")
	fmt.Scanln(&input)
}
