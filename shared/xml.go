package shared

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
