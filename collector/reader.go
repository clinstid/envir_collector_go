package main

import (
	"encoding/xml"
	"log"

	"github.com/clinstid/envir_collector_go/shared"
)

func readValues(client *shared.EnvirClient, c chan shared.XMLMessage) {
	client.Open()

	defer client.Close()

	log.Println("[r]: Waiting for data...")

	msg := make([]byte, 0)
	buf := make([]byte, 1)
	for {
		if _, err := client.Read(buf); err != nil {
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
