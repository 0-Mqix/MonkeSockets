# MonkeSockets
websockets with events and rooms made with "github.com/gorilla/websocket"

## Use the Client
####  Java Script
```js
const monkeSocket = new MonkeSocket("/ws")

monkeSocket.on("message:", (message) => {
    //this runs when there is a message: event
})

monkeSocket.onOpen(() => {
    //this runs when the connection opens
})

monkeSocket.onClose(() => {
    //this runs when the connection closes
})

//example for sending an event + its message
monkeSocket.send("message:", "hi server")
```

## Use the Server
#### Go
```go
package main

import (
	"github.com/0-Mqix/MonkeSockets"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)


func main() {
	r := fiber.New()

	r.Static("/", "./public")

	rooms := make(map[string]*MonkeSockets.Room)
	room := MonkeSockets.New()
	go room.Run()

	rooms["default"] = room

	r.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	r.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		client := &MonkeSockets.Client{Rooms: rooms, Conn: c, Channel: make(chan []byte), Closed: false}

		defer func() {
			client.Disconnect()
		}()

		for name, room := range rooms {
			room.Register <- client
			client.Rooms[name] = room
		}

		go client.Writer()
		client.Reader()
	}))

	r.Listen(":80")
}
```