package gocc

import (
	"reflect"
	"strings"
	"testing"
)

var testData = `<msg><src>CC128-v0.11</src><dsb>00089</dsb><time>13:02:39</time><tmpr>18.7</tmpr><sensor>1</sensor><id>01234</id><type>1</type><ch1><watts>00345</watts></ch1><ch2><watts>02151</watts></ch2><ch3><watts>00000</watts></ch3></msg>
`

func newInt(v int) *int {
	return &v
}

func newFloat32(v float32) *float32 {
	return &v
}

func TestReadSensorUpdate(t *testing.T) {
	r := strings.NewReader(testData)
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
