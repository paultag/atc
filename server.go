package main

import (
	"fmt"
	"time"

	"pault.ag/go/atc/atc"
)

type Location struct {
	Latitude  string
	Longitude string
}

type Aircraft struct {
	HexIdent    string
	Altitude    string
	Callsign    string
	Squawk      string
	GroundSpeed string
	Track       string
	Locations   []Location

	LastUpdated time.Time
}

func ohshit(err error) {
	if err != nil {
		panic(err)
	}
}

func Update(aircraft *Aircraft, message *atc.Message) {
	switch message.TransmissionType {
	case "1":
		aircraft.Callsign = message.Callsign
	case "2":
		aircraft.Altitude = message.Altitude
		aircraft.GroundSpeed = message.GroundSpeed
		aircraft.Track = message.Track
		aircraft.Locations = append(aircraft.Locations, Location{
			Latitude:  message.Latitude,
			Longitude: message.Longitude,
		})
	case "3":
		aircraft.Altitude = message.Altitude
		aircraft.Locations = append(aircraft.Locations, Location{
			Latitude:  message.Latitude,
			Longitude: message.Longitude,
		})
	case "4":
		aircraft.GroundSpeed = message.GroundSpeed
	case "5":
		aircraft.Altitude = message.Altitude
	case "6":
		aircraft.Squawk = message.Squawk
	case "7":
		aircraft.Altitude = message.Altitude
	}
}

func main() {
	data := make(chan *atc.Message, 100)
	go func() {
		defer panic("Stream died")
		// atc.Stream("aircraft.paultag.house:30006", data)
		atc.Stream("localhost:30006", data)
	}()

	airspace := map[string]*Aircraft{}
	go func() {
		for message := range data {
			var ok bool
			var aircraft *Aircraft

			if aircraft, ok = airspace[message.HexIdent]; !ok {
				aircraft = &Aircraft{}
				airspace[message.HexIdent] = aircraft
			}
			Update(aircraft, message)
		}
	}()

	for {
		for _, aircraft := range airspace {
			fmt.Printf("%s\n", aircraft)
		}
		fmt.Printf("\n\n\n\n\n\n")
		time.Sleep(time.Second * 1)
	}
}
