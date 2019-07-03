# socketchan

> socketchan, burning through WebSockets here alone

socketchan is a thin wrapper over [Gorilla WebSockets](https://github.com/gorilla/websocket) which simply exposes a WebSocket connection as 2 Go channels.

See [`socketchan/cmd/lurk`](./cmd/lurk/main.go) for an example program that logs all messages received ([`hark`](./cmd/hark/main.go) will send the time every 5 seconds):

```go
package main

import (
	"log"
	"net/url"

	sc "github.com/notduncansmith/socketchan"
)

func main() {
	u, _ := url.Parse("ws://localhost:8080/room/demo")
	client := sc.NewClient(*u, 1024) // channels will buffer up to 1024 messages

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
```

Copyright © 2019 Duncan Smith