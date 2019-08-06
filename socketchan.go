package socketchan

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client is a Websocket client that connects to an endpoint
type Client struct {
	mut        *sync.RWMutex
	connected  bool
	bufferSize int
	Conn       *websocket.Conn
	Endpoint   url.URL
	Incoming   chan []byte
	Outgoing   chan []byte
}

// NewClient returns a new Client whose channels will buffer up to `bufferSize` messages
func NewClient(endpoint url.URL, bufferSize int) *Client {
	incoming := make(chan []byte, bufferSize)
	outgoing := make(chan []byte, bufferSize)
	mut := sync.RWMutex{}
	client := Client{&mut, false, bufferSize, nil, endpoint, incoming, outgoing}
	return &client
}

// Connect opens the connection to the WebSocket endpoint and starts goroutines to populate the channels
func (c *Client) Connect() error {
	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Duration(45) * time.Second,
		Proxy:            http.ProxyFromEnvironment,
	}

	conn, _, err := dialer.Dial(c.Endpoint.String(), nil)

	if err != nil {
		return err
	}

	c.Conn = conn
	c.doWithLock(func() {
		c.connected = true
		c.Incoming = make(chan []byte, c.bufferSize)
		c.Outgoing = make(chan []byte, c.bufferSize)
	})

	go func() {
		defer c.Close()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if strings.Contains(err.Error(), "1006") {
					// socket closed
					return
				}
				log.Println("socketchan read err:", err)
				return
			}
			c.Incoming <- message
		}
	}()

	go func() {
		defer c.Close()

		for out := range c.Outgoing {
			err := conn.WriteMessage(websocket.TextMessage, out)
			if err != nil {
				log.Println("socketchan write err:", err)
				return
			}
		}
	}()

	return nil
}

// Close closes the WebSocket connection, which fires the OnClose callback
func (c *Client) Close() error {
	c.doWithLock(func() {
		if c.connected {
			c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			close(c.Incoming)
			close(c.Outgoing)
			c.connected = false
		}
	})

	return nil
}

// Connected returns whether or not the client is connected
func (c *Client) Connected() bool {
	return c.withRLock(func() interface{} {
		return c.connected
	}).(bool)
}

func (c *Client) doWithLock(f func()) {
	c.mut.Lock()
	defer c.mut.Unlock()
	f()
}

func (c *Client) withRLock(f func() interface{}) interface{} {
	c.mut.RLock()
	defer c.mut.RUnlock()
	return f()
}
