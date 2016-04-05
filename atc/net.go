package atc

import (
	"bufio"
	"log"
	"net"
)

func Stream(host string, output chan<- *Message) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}
	reader := bufio.NewReader(conn)

	for {
		el, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("%s\n", err)
			return
		}

		msg, err := Parse(el)
		if err != nil {
			log.Printf("%s\n", err)
			return
		}

		output <- msg
	}
}
