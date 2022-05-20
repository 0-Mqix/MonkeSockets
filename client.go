package MonkeSockets

import (
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Rooms   map[string]*Room
	Conn    *websocket.Conn
	Channel chan []byte
	Closed  bool
}

type SocketMessage struct {
	Client  *Client
	Message []byte
}

func (c *Client) Disconnect() {
	if c.Closed {
		return
	}

	c.Closed = true
	close(c.Channel)

	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})

	for _, r := range c.Rooms {
		r.unregister <- c
	}
}

func (c *Client) Writer() {
	for {
		message, ok := <-c.Channel

		if !ok {
			c.Disconnect()
			return
		}

		c.Conn.WriteMessage(websocket.TextMessage, message)
	}

}

func (c *Client) Reader() {
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			c.Disconnect()
			return
		}

		for _, r := range c.Rooms {
			r.message <- SocketMessage{Message: message, Client: c}
		}
	}
}

//sends message to the client
func (c *Client) Send(event string, message []byte) {
	select {
	case c.Channel <- append([]byte(event), message...):
	default:
		c.Disconnect()
	}
}
