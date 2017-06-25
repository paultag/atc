package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dustin/go-broadcast"
	"golang.org/x/net/websocket"
	"pault.ag/go/atc/atc"
)

var lock = sync.RWMutex{}

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
	Location    Location

	LastUpdated time.Time
}

func (a *Aircraft) UpdateLoaction(location Location) {
	a.Location = location
	a.Locations = append(a.Locations, location)
}

func Update(aircraft *Aircraft, message *atc.Message) {
	aircraft.LastUpdated = time.Now()
	switch message.TransmissionType {
	case "1":
		aircraft.Callsign = message.Callsign
	case "2":
		aircraft.Altitude = message.Altitude
		aircraft.GroundSpeed = message.GroundSpeed
		aircraft.Track = message.Track
		aircraft.UpdateLoaction(Location{
			Latitude:  message.Latitude,
			Longitude: message.Longitude,
		})
	case "3":
		aircraft.Altitude = message.Altitude
		aircraft.UpdateLoaction(Location{
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

type Airspace map[string]*Aircraft

func (airspace Airspace) All() []*Aircraft {
	ret := []*Aircraft{}
	lock.RLock()
	defer lock.RUnlock()
	for _, el := range airspace {
		ret = append(ret, el)
	}
	return ret
}

func (airspace Airspace) Get(id string) *Aircraft {
	lock.RLock()
	defer lock.RUnlock()
	var ok bool
	var aircraft *Aircraft
	if aircraft, ok = airspace[id]; !ok {
		aircraft = &Aircraft{HexIdent: id}
		airspace[id] = aircraft
	}
	return aircraft
}

func (a Airspace) Clean() []string {
	removed := []string{}
	for id, el := range a {
		delta := time.Now().Sub(el.LastUpdated)
		if delta > time.Second*2400 {
			a.Delete(id)
			removed = append(removed, id)
		}
	}
	return removed
}

func (a Airspace) Delete(id string) {
	lock.RLock()
	defer lock.RUnlock()
	delete(a, id)
}

type Message struct {
	Location Location
	Track    []Location
	HexIdent string
}

func main() {
	data := make(chan *atc.Message, 100)
	go func() {
		defer panic("Stream died")
		atc.Stream("aircraft.paultag.house:30006", data)
		// atc.Stream("localhost:30006", data)
	}()

	clients := broadcast.NewBroadcaster(100)
	defer clients.Close()
	airspace := Airspace{}

	go func() {
		for message := range data {
			aircraft := airspace.Get(message.HexIdent)
			Update(aircraft, message)
			if len(aircraft.Locations) == 0 {
				continue
			}
			clients.Submit(*aircraft)
		}
	}()

	go func() {
		for {
			removed := airspace.Clean()
			for _, id := range removed {
				fmt.Printf("Cleaned out ID: %s\n", id)
				clients.Submit(Aircraft{HexIdent: id})
			}
			time.Sleep(time.Second * 30)
		}
	}()

	http.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		updates := make(chan interface{}, 100)
		clients.Register(updates)
		defer clients.Unregister(updates)
		output := json.NewEncoder(ws)

		for _, el := range airspace.All() {
			output.Encode(*el)
		}

		for el := range updates {
			aircraft := el.(Aircraft)
			output.Encode(aircraft)
		}
	}))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
