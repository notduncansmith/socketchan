package main

import (
	"log"
	"net/url"
	"time"

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

	ticks := time.NewTicker(5 * time.Second)

	for t := range ticks.C {
		msg := "{\"time\": \"" + t.String() + "\"}"
		log.Println("Sending msg: " + msg)
		client.Outgoing <- []byte(msg)
	}
}
