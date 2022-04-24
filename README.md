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
func (room *MonkeSockets.Room) WebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return nil
	}
	client := &Client{room: MonkeSockets.room, conn: MonkeSockets.conn, send: make(chan []byte, 256)}
	client.room.register <- client

	go client.WritePump()
	go client.ReadPump()
	return nil
}

func main() {
	e := echo.New()
    e.Static("/static/js/MonkeSocket.js")

    r := MonkeSockets.New()
    r.Run()

    e.("/ws", r.WebSocket)
	e.Logger.Fatal(e.Start(":8080"))
}
```