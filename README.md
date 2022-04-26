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
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Connect(room *MonkeSockets.Room, c echo.Context) {
	conn, err := Upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return
	}

	client := &MonkeSockets.Client{Rooms: make([]*MonkeSockets.Room, 0), Conn: conn, Channel: make(chan []byte, 256), Echo: c}
	client.Rooms = append(client.Rooms, room)

	room.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func main() {
	e := echo.New()
	e.File("/static/js/index.js", "index.js")
	e.File("/", "index.html")

	r := MonkeSockets.New()
	go r.Run()

	e.GET("/ws", func(c echo.Context) error {
		Connect(r, c)
		return nil
	})

	e.Logger.Fatal(e.Start(":8080"))
}
```