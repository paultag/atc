package main

import (
	"fmt"

	"pault.ag/go/atc/atc"
)

func ohshit(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	data := make(chan *atc.Message, 100)
	go func() {
		defer panic("Stream died")
		atc.Stream("localhost:30006", data)
	}()

	for message := range data {
		when, err := message.GeneratedAt()
		ohshit(err)
		fmt.Printf("%s\n", when)
	}
}
