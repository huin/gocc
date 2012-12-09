// gocc is an unofficial reader for [Current Cost](http://www.currentcost.com/)
// energy monitoring system data for the [Go](http://golang.org) (Golang)
// programming language. This software is NOT endorsed by or affiliated with
// Current Cost.
//
// Typical usage:
//
//   msgReader, err := gocc.NewSerialMessageReader("/dev/ttyUSB0")
//   if err != nil {
//     // deal with err
//   }
//   for {
//     if msg, err := msgReader.ReadMessage(); err != nil {
//       // deal with err
//     } else {
//       // use msg
//     }
//   }
package gocc

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/huin/goserial"
)

var (
	ErrLineTooLong = errors.New("line too long")
	ErrLineEmpty   = errors.New("line empty")
)

type MessageReader struct {
	// Closer is an optional field, that will have its Close method called when
	// MessageReader.Close is called.
	Closer io.Closer

	// Reader reads lines received from a Current Cost unit. The buffer size must
	// be large enough to hold a single line from the Current Cost unit.
	Reader *bufio.Reader
}

// NewMessageReader creates a new MessageReader that will read lines from. If r
// implements io.Closer, then it will be closed when MessageReader.Close is
// called.
func NewMessageReader(r io.Reader) *MessageReader {
	c, _ := r.(io.Closer)
	return &MessageReader{
		Closer: c,
		Reader: bufio.NewReaderSize(r, 16*1024),
	}
}

// NewSerialMessageReader opens the named serial port, configures it for
// reading Current Cost data, and returns a MessageReader for doing so.
func NewSerialMessageReader(serialPath string) (*MessageReader, error) {
	serialConfig := &goserial.Config{
		Name:     serialPath,
		Baud:     57600,
		Parity:   goserial.ParityNone,
		Size:     goserial.Byte8,
		StopBits: goserial.StopBits1,
	}

	var serial io.ReadCloser
	var err error
	if serial, err = goserial.OpenPort(serialConfig); err != nil {
		return nil, err
	}

	return NewMessageReader(serial), nil
}

func (reader *MessageReader) String() string {
	return fmt.Sprintf("<Current Cost scraper from %s>", reader.Reader)
}

// Close closes the underlying Closer (if any is set).
func (reader *MessageReader) Close() error {
	if reader.Closer != nil {
		return reader.Closer.Close()
	}
	return nil
}

func (reader *MessageReader) ReadMessage() (*Message, error) {
	line, isPrefix, err := reader.Reader.ReadLine()
	if isPrefix {
		return nil, ErrLineTooLong
	} else if err != nil {
		return nil, err
	}

	// The Current Cost unit seems to occasionally insert a \xfc byte at the
	// start of a line. Discard if present.
	if len(line) > 0 && line[0] == 0xfc {
		line = line[1:]
	}

	if len(line) == 0 {
		return nil, ErrLineEmpty
	}

	msg := new(Message)
	if err = xml.Unmarshal(line, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

type SensorType int

const (
	SensorElectricity = SensorType(1)
)

type UnitsType string

const (
	UnitKWHr = UnitsType("kwhr")
)

// Message is the top-level data type representing data from a Current Cost unit.
type Message struct {
	// Always present fields:
	Src            string `xml:"src"`
	DaysSinceBirth int    `xml:"dsb"`
	TimeOfDay      string `xml:"time"`

	// Present in real-time updates:
	Temperature *float32    `xml:"tmpr"`
	Sensor      *int        `xml:"sensor"`
	ID          *int        `xml:"id"`
	Type        *SensorType `xml:"type"`
	Channel1    *Channel    `xml:"ch1"`
	Channel2    *Channel    `xml:"ch2"`
	Channel3    *Channel    `xml:"ch3"`

	// Present in history updates:
	History *History `xml:"hist"`
}

type Channel struct {
	Watts int `xml:"watts"`
}

type History struct {
	DaysSinceWipe int           `xml:"dsw"`
	Type          SensorType    `xml:"type"`
	Units         UnitsType     `xml:"units"`
	Data          []HistoryData `xml:"data"`
}

type HistoryData struct {
	Sensor int `xml:"sensor"`

	// Sometimes present:
	Units *UnitsType `xml:"units"`

	// Values over time.
	Values []HistoryDataPoint `xml:",any"`
}

type HistoryDataPoint struct {
	XMLName xml.Name // Represents time range (e.g "h024" meaning 22 to 24 hours ago).
	Value   float32  `xml:",chardata"`
}
