package main

import (
	"log"

	"pault.ag/go/atc/atc"
)

/*
 * '1': callsign
 * '2': alt, GS, trk, lat, long, grnd
 * '3': alt, lat, long, alert, emerg, spi, ground
 * '4': ground_speed, track
 * '5': alt, alert, spi, grnd
 * '6': squawk
 * '7': altitude
 */

func main() {
	data := make(chan *atc.Message, 100)
	go func() {
		defer panic("Stream died")
		atc.Stream("192.168.1.117:30006", data)
	}()

	for el := range data {
		switch el.TransmissionType {
		case "3", "2":
			log.Printf("%s %s %s\n", el.HexIdent, el.Latitude, el.Longitude)
		}
	}
}
