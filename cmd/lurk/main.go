package main

import (
	"flag"
	"log"
	"net/url"

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
