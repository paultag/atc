package main

import (
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
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
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "flights",
		Password: "flights",
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	data := make(chan *atc.Message, 100)
	go func() {
		defer panic("Stream died")
		atc.Stream("192.168.1.117:30006", data)
	}()

	for el := range data {
		fields := map[string]interface{}{}
		tags := map[string]string{"transponder": el.HexIdent}
		collection := ""

		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  "flights",
			Precision: "us",
		})
		if err != nil {
			log.Fatalln("Error: ", err)
		}

		switch el.TransmissionType {
		case "3", "2":
			log.Printf("%s %s %s\n", el.HexIdent, el.Latitude, el.Longitude)
			collection = "locations"
			fields = map[string]interface{}{
				"altitude":  el.Altitude,
				"latitude":  el.Latitude,
				"longitude": el.Longitude,
			}
		case "1":
			log.Printf("%s %s\n", el.HexIdent, el.Callsign)
			collection = "callsigns"
			fields = map[string]interface{}{
				"callsign": el.Callsign,
			}
		case "6":
			log.Printf("%s %s\n", el.HexIdent, el.Squawk)
			collection = "squawks"
			fields = map[string]interface{}{
				"squawk": el.Squawk,
			}

		default:
			continue
		}

		if collection == "" {
			panic("No collection")
			continue
		}

		pt, err := client.NewPoint(collection, tags, fields, time.Now())
		if err != nil {
			log.Fatalln("Error: ", err)
		}
		bp.AddPoint(pt)
		c.Write(bp)
	}
}
