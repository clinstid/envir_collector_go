package shared

import (
	"errors"
	"log"
	"os"

	serial "github.com/mikepb/go-serial"
)

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

type EnvirClient struct {
	port      *serial.Port
	usbDevice string
	options   serial.Options
}

func NewEnvirClient(bitRate, dataBits, stopBits int, usbDevice string) (*EnvirClient, error) {
	c := new(EnvirClient)

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
				c.usbDevice = e
				break
			}
		}
	} else {
		c.usbDevice = usbDevice
	}

	if c.usbDevice == "" {
		log.Panic("Couldn't find a usable USB device!")
		return nil, errors.New("Couldn't find a usable USB device!")
	}
	c.options = serial.Options{BitRate: bitRate, DataBits: dataBits, StopBits: stopBits}

	return c, nil
}

func (c *EnvirClient) Open() error {
	var err error
	c.port, err = c.options.Open(c.usbDevice)
	if err != nil {
		log.Panic("Failed to open usb device", c.usbDevice, ":", err)
		return err
	}
	return nil
}

func (c *EnvirClient) Close() error {
	return c.port.Close()
}

func (c *EnvirClient) Read(b []byte) (int, error) {
	return c.port.Read(b)
}
