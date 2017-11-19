package shared

import (
	"bytes"
	"database/sql"
	"fmt"
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

type DBClient struct {
	dbOpenString string
}

func NewDBClient(dbHost, dbUser, dbPassword, dbName string) *DBClient {
	c := new(DBClient)
	c.dbOpenString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s", dbHost, dbUser, dbPassword, dbName)
	return c
}

func (c *DBClient) open() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.dbOpenString)
	if err != nil {
		log.Panic("Failed to connect to database")
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Panic("Failed to ping db:", err)
		return nil, err
	}
	return db, err
}

func (c *DBClient) WriteMessage(x XMLMessage) error {
	log.Printf("[w]: shared.XMLMessage object: %+v\n", x)
	dbMessage := buildDBMessage(x)
	insertCmd := genInsert(dbMessage)
	log.Printf("[w]: Insert command: %+v\n", insertCmd)
	db, err := c.open()
	if err != nil {
		return err
	}
	res, err := db.Exec(insertCmd)
	if err != nil {
		log.Panic("Failed INSERT:", err)
	}
	log.Printf("[w]: Insert result: %+v\n", res)
	err = db.Close()
	if err != nil {
		log.Panic("Failed to close the database connection:", err)
	}
	return err
}
