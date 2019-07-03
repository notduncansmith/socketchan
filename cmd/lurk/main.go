package main

import (
	"log"
	"net/url"

	sc "github.com/notduncansmith/socketchan"
)

func main() {
	u, _ := url.Parse("ws://localhost:8080/room/demo")
	client := sc.NewClient(*u, 1024)

	err := client.Connect()

	if err != nil {
		log.Fatalln("dial err", err)
	}

	log.Println("Connected to " + u.String())

	for {
		select {
		case rawMsg, stillOpen := <-client.Incoming:
			if !stillOpen {
				log.Println("Socket closed, exiting")
				return
			}
			log.Println("\n--\n" + string(rawMsg) + "\n--\n")
		}
	}
}
