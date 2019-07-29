package socketchan

import (
	"log"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

// Client is a Websocket client that connects to an endpoint
type Client struct {
	Conn     *websocket.Conn
	Endpoint url.URL
	Incoming chan []byte
	Outgoing chan []byte
}

// NewClient returns a new Client whose channels will buffer up to `bufferSize` messages
func NewClient(endpoint url.URL, bufferSize int) *Client {
	incoming := make(chan []byte, bufferSize)
	outgoing := make(chan []byte, bufferSize)
	client := Client{nil, endpoint, incoming, outgoing}
	return &client
}

// Connect opens the connection to the WebSocket endpoint and starts goroutines to populate the channels
func (c *Client) Connect() error {
	e := c.Endpoint.String()

	conn, _, err := websocket.DefaultDialer.Dial(e, nil)

	if err != nil {
		return err
	}

	c.Conn = conn

	go func() {
		defer close(c.Incoming)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if strings.Contains(err.Error(), "1006") {
					// socket closed
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

		for out := range c.Outgoing {
			err := conn.WriteMessage(websocket.TextMessage, out)
			if err != nil {
				log.Println("write err:", err)
				return
			}
		}
	}()

	return nil
}

// Close closes the WebSocket connection, which fires the OnClose callback
func (c *Client) Close() error {
	err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	close(c.Outgoing)
	return nil
}
