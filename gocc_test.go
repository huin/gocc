package gocc

import (
	"encoding/xml"
	"reflect"
	"strings"
	"testing"
)

// Example taken from http://www.currentcost.com/cc128/xml.htm
const testRealtimeData = `<msg><src>CC128-v0.11</src><dsb>00089</dsb><time>13:02:39</time><tmpr>18.7</tmpr><sensor>1</sensor><id>01234</id><type>1</type><ch1><watts>00345</watts></ch1><ch2><watts>02151</watts></ch2><ch3><watts>00000</watts></ch3></msg>
`

func newInt(v int) *int {
	return &v
}

func newFloat32(v float32) *float32 {
	return &v
}

func newUnits(v UnitsType) *UnitsType {
	return &v
}

func xmlName(v string) xml.Name {
	return xml.Name{Local: v}
}

func TestReadRealtimeSensorUpdate(t *testing.T) {
	r := strings.NewReader(testRealtimeData)
	msgReader := NewMessageReader(r)

	msg, err := msgReader.ReadMessage()
	expected := &Message{
		Src:            "CC128-v0.11",
		DaysSinceBirth: 89,
		TimeOfDay:      "13:02:39",

		Temperature: newFloat32(18.7),
		Sensor:      newInt(1),
		ID:          newInt(1234),

		Type:     (*SensorType)(newInt(1)),
		Channel1: &Channel{Watts: 345},
		Channel2: &Channel{Watts: 2151},
		Channel3: &Channel{Watts: 0},
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if !reflect.DeepEqual(expected, msg) {
		t.Errorf("mismatched expected result\nexpected %#v\n     got %#v", expected, msg)
	}
}

// Example taken from http://www.currentcost.com/cc128/xml.htm
const testTwoHourlyHistoryData = `<msg><src>CC128-v0.11</src><dsb>00089</dsb><time>13:10:50</time><hist><dsw>00032</dsw><type>1</type><units>kwhr</units><data><sensor>0</sensor><h024>001.1</h024><h022>000.9</h022><h020>000.3</h020><h018>000.4</h018></data><data><sensor>9</sensor><units>kwhr</units><h024>000.0</h024><h022>000.0</h022><h020>000.0</h020><h018>000.0</h018></data></hist></msg>`

func TestTwoHourlyHistory(t *testing.T) {
	r := strings.NewReader(testTwoHourlyHistoryData)
	msgReader := NewMessageReader(r)

	msg, err := msgReader.ReadMessage()
	expected := &Message{
		Src:            "CC128-v0.11",
		DaysSinceBirth: 89,
		TimeOfDay:      "13:10:50",
		History: &History{
			DaysSinceWipe: 32,
			Type:          SensorElectricity,
			Units:         UnitKWHr,
			Data: []HistoryData{
				{
					Sensor: 0,
					Values: []HistoryDataPoint{
						{XMLName: xmlName("h024"), Value: 1.1},
						{XMLName: xmlName("h022"), Value: 0.9},
						{XMLName: xmlName("h020"), Value: 0.3},
						{XMLName: xmlName("h018"), Value: 0.4},
					},
				},
				{
					Sensor: 9,
					Units:  newUnits(UnitKWHr),
					Values: []HistoryDataPoint{
						{XMLName: xmlName("h024"), Value: 0},
						{XMLName: xmlName("h022"), Value: 0},
						{XMLName: xmlName("h020"), Value: 0},
						{XMLName: xmlName("h018"), Value: 0},
					},
				},
			},
		},
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if !reflect.DeepEqual(expected, msg) {
		t.Errorf("mismatched expected result\nexpected %#v\n     got %#v", expected, msg)
	}
}
