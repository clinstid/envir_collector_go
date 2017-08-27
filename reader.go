package main

import (
	"encoding/xml"
	"log"

	"github.com/mikepb/go-serial"
)

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
