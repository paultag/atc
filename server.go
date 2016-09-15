package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"

	"golang.org/x/net/websocket"
	"pault.ag/go/atc/atc"
)

func main() {
	data := make(chan *atc.Message, 100)
	go func() {
		defer panic("Stream died")
		atc.Stream("ladon.paultag.house:30006", data)
	}()

	clients := NewClients()
	go func() {
		defer panic("Lost ATC feed")
		for el := range data {
			switch el.TransmissionType {
			case "3", "2":
				data, err := json.Marshal(map[string]interface{}{
					"hex":       el.HexIdent,
					"altitude":  el.Altitude,
					"latitude":  el.Latitude,
					"longitude": el.Longitude,
				})
				if err != nil {
					log.Printf("%s\n", err)
					continue
				}
				clients.Broadcast(data)
			default:
				continue
			}
		}
	}()

	http.Handle("/atc", websocket.Handler(func(ws *websocket.Conn) {
		client := NewClient(ws)
		clients.Register(client)
		io.Copy(ioutil.Discard, ws)
	}))

	err := http.ListenAndServe(":1090", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
