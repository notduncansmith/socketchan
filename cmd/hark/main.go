package main

import (
	"flag"
	"log"
	"net/url"
	"time"

	sc "github.com/notduncansmith/socketchan"
)

func main() {
	flag.Parse()
	var addr = flag.Arg(0)
	u, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	println("Connecting to " + u.String())
	client := sc.NewClient(*u, 1024)

	err = client.Connect()

	ticks := time.NewTicker(5 * time.Second)

	for t := range ticks.C {
		msg := "{\"time\": \"" + t.String() + "\"}"
		log.Println("Sending msg: " + msg)
		client.Outgoing <- []byte(msg)
	}
}
