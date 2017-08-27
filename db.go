package main

import (
	"bytes"
	"log"
	"text/template"
	"time"

	"github.com/lib/pq"
)

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
