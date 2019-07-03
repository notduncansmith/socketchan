package socketchan

import (
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

// Client is a Websocket client that connects to an endpoint
type Client struct {
	Endpoint url.URL
	Incoming chan []byte
	Outgoing chan []byte
	Close    chan struct{}
}

// NewClient returns a new client whose channels will buffer up to `bufferSize` messages
func NewClient(endpoint url.URL, bufferSize int) *Client {
	incoming := make(chan []byte, bufferSize)
	outgoing := make(chan []byte, bufferSize)
	close := make(chan struct{})
	client := Client{endpoint, incoming, outgoing, close}
	return &client
}

// Connect opens the connection to the WebSocket endpoint and starts goroutines to populate the channels
func (c *Client) Connect() error {
	e := c.Endpoint.String()

	conn, _, err := websocket.DefaultDialer.Dial(e, nil)

	if err != nil {
		return err
	}

	go func() {
		defer close(c.Close)
		defer close(c.Incoming)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if strings.Contains(err.Error(), "1006") {
					log.Println("socket closed")
					return
				}
				log.Println("read err:", err)
				return
			}
			c.Incoming <- message
		}
	}()

	go func() {
		defer conn.Close()

		for {
			select {
			case out := <-c.Outgoing:
				err := conn.WriteMessage(websocket.TextMessage, out)
				if err != nil {
					log.Println("write err:", err)
					return
				}
			case _, stillOpen := <-c.Close:
				if !stillOpen {
					err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Println("write close err:", err)
					}
					return
				}
			}
		}
	}()

	return nil
}
