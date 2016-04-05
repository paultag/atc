package main

import (
	"log"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

func main() {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "flights",
		Password: "flights",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "flights",
		Precision: "us",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	// Create a point and add to batch
	tags := map[string]string{"transponder": "ABCDAEF"}
	fields := map[string]interface{}{
		"altitude":  10.1,
		"latitude":  53.3,
		"longitude": 46.6,
	}
	pt, err := client.NewPoint("blips", tags, fields, time.Now())

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(pt)

	// Write the batch
	c.Write(bp)
}
