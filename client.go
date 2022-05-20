package MonkeSockets

import (
	"fmt"
	"log"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Rooms   map[string]*Room
	Conn    *websocket.Conn
	Channel chan []byte
}

type SocketMessage struct {
	Client  *Client
	Message []byte
}

func (c *Client) Disconnect() {
	fmt.Println("disconect")
	c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
	close(c.Channel)
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
			log.Fatal(err)
			c.Disconnect()
		}

		log.Printf("incomming: %s", message)
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
