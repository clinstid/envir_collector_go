package main

import (
	"encoding/xml"
	"log"

	"github.com/clinstid/envir_collector_go/shared"

	"github.com/mikepb/go-serial"
)

func readValues(bitRate, dataBits, stopBits int, usbDevice string, c chan shared.XMLMessage) {
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
		var exampleMessage shared.XMLMessage
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
				log.Println("[r]: shared.XMLMessage string:", msgStr)

				var message shared.XMLMessage
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
