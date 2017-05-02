package main

import (
	"bytes"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/lib/pq"
	"github.com/mikepb/go-serial"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"
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

// XMLMessage XML element
type XMLMessage struct {
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

// DBMessage database message format
type DBMessage struct {
	Src      string
	Dsb      string
	Time     string
	TmprF    float32
	Sensor   int
	DeviceID string
	Ch1Watts int
	Ch2Watts int
	Ch3Watts int
}

func buildDBMessage(x XMLMessage) DBMessage {
	dbMessage := DBMessage{
		Src:      x.Src,
		Dsb:      x.Dsb,
		Time:     string(pq.FormatTimestamp(time.Now())),
		TmprF:    x.TmprF,
		Sensor:   x.Sensor,
		DeviceID: x.ID,
		Ch1Watts: x.Ch1.Watts,
		Ch2Watts: x.Ch2.Watts,
		Ch3Watts: x.Ch3.Watts,
	}
	return dbMessage
}

func readValues(bitRate, dataBits, stopBits int, usbDevice string, c chan XMLMessage) {
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
		var exampleMessage XMLMessage
		err := xml.Unmarshal([]byte(exampleXml), &exampleMessage)
		if err != nil {
			log.Panic("Failed on example xml:", err)
		}
		log.Println("Example XML result:", exampleMessage)
	*/
	options := serial.Options{BitRate: bitRate, DataBits: dataBits, StopBits: stopBits}
	log.Println("[r]: Connecting to device", usbDevice)
	p, err := options.Open(usbDevice)
	if err != nil {
		log.Panic("Failed to open usb device:", err)
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
				log.Println("[r]: XMLMessage string:", msgStr)

				var message XMLMessage
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

func genInsert(m DBMessage) string {
	const templateString string = "INSERT INTO energydash " +
		"(src, dsb, time, tmprf, sensor, device_id, ch1_watts, ch2_watts, ch3_watts) " +
		"VALUES ('{{.Src}}', '{{.Dsb}}', '{{.Time}}', '{{.TmprF}}', '{{.Sensor}}', " +
		"'{{.DeviceID}}', '{{.Ch1Watts}}', '{{.Ch2Watts}}', '{{.Ch3Watts}}')"

	tmpl, err := template.New("insert").Parse(templateString)
	if err != nil {
		log.Panic("Failed to parse template:", err)
	}
	var insertCmdBytes bytes.Buffer
	err = tmpl.Execute(&insertCmdBytes, m)
	if err != nil {
		log.Panic("Failed to generate insert command:", err)
	}
	return insertCmdBytes.String()
}

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

func main() {
	bitRate := getenvAsInt("ENVIR_SERIAL_BIT_RATE", 57600)
	dataBits := getenvAsInt("ENVIR_SERIAL_DATA_BITS", 8)
	stopBits := getenvAsInt("ENVIR_SERIAL_STOP_BITS", 1)
	usbDevice := getenv("ENVIR_SERIAL_USB_DEVICE", "/dev/ttyUSB0")

	// Potential USB device files
	usbDeviceList := []string{
		"/dev/ttyUSB0",
		"/dev/ttyUSB1",
		"/dev/ttyUSB2",
		"/dev/ttyUSB3",
	}

	// Attempt to check if the device exists by going through usbDevice.
	// Hopefully the specified device will work.
	if _, err := os.Stat(usbDevice); os.IsNotExist(err) {
		for _, e := range usbDeviceList {
			if _, err := os.Stat(e); err == nil {
				usbDevice = e
				break
			}
		}
	}

	dbHost := getenv("ENVIR_DB_HOST", "yoda") // TODO: Switch to localhost
	dbUser := getenv("ENVIR_DB_USER", "energydash")
	dbPassword := getenv("ENVIR_DB_PASSWORD", "energydash")
	dbName := getenv("ENVIR_DB_NAME", "energydash")

	var c = make(chan XMLMessage, 1000)

	go readValues(bitRate, dataBits, stopBits, usbDevice, c)
	go writeValues(c, dbHost, dbUser, dbPassword, dbName)

	var input string
	fmt.Printf("Press enter to stop execution...")
	fmt.Scanln(&input)
}
