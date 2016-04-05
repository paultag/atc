package atc

import (
	"fmt"
	"strings"
)

//
//
type Message struct {
	MessageType      string
	TransmissionType string
	SessionID        string
	AircraftID       string
	HexIdent         string
	FlightID         string
	DateGenerated    string
	TimeGenerated    string
	DateLogged       string
	TimeLogged       string
	Callsign         string
	Altitude         string
	GroundSpeed      string
	Track            string
	Latitude         string
	Longitude        string
	VerticalRate     string
	Squawk           string
	Alert            string
	Emergency        string
	SPI              string
	OnGround         string
}

func (m *Message) Marshal() []byte {
	return []byte(
		strings.Join([]string{
			m.MessageType,
			m.TransmissionType,
			m.SessionID,
			m.AircraftID,
			m.HexIdent,
			m.FlightID,
			m.DateGenerated,
			m.TimeGenerated,
			m.DateLogged,
			m.TimeLogged,
			m.Callsign,
			m.Altitude,
			m.GroundSpeed,
			m.Track,
			m.Latitude,
			m.Longitude,
			m.VerticalRate,
			m.Squawk,
			m.Alert,
			m.Emergency,
			m.SPI,
			m.OnGround,
		}, ","),
	)
}

func (m *Message) Unmarshal(data []byte) error {
	els := strings.Split(string(data), ",")
	if len(els) != 22 {
		return fmt.Errorf("ErrorBad length%d", len(els))
	}

	m.MessageType = els[0]
	m.TransmissionType = els[1]
	m.SessionID = els[2]
	m.AircraftID = els[3]
	m.HexIdent = els[4]
	m.FlightID = els[5]
	m.DateGenerated = els[6]
	m.TimeGenerated = els[7]
	m.DateLogged = els[8]
	m.TimeLogged = els[9]
	m.Callsign = els[10]
	m.Altitude = els[11]
	m.GroundSpeed = els[12]
	m.Track = els[13]
	m.Latitude = els[14]
	m.Longitude = els[15]
	m.VerticalRate = els[16]
	m.Squawk = els[17]
	m.Alert = els[18]
	m.Emergency = els[19]
	m.SPI = els[20]
	m.OnGround = els[21]

	return nil
}

//

func Parse(data string) (*Message, error) {
	ret := Message{}
	return &ret, ret.Unmarshal([]byte(data))
}
