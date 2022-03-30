package main

import (
	"projects/chat/MonkeSockets"

	"github.com/labstack/echo/v4"
)

func home(e echo.Context) error {
	return e.HTML(200, "")
}

type SocketMessageJsonIn struct {
	Msg string `json:"msg"`
}

func main() {
	e := echo.New()
	e.File("/", "home.html")
	e.File("/test", "test.html")
	e.Static("/static", "./MonkeSocket")
	//SOCKETS
	room := MonkeSockets.NewRoom()

	// room.Events["chat:"] = func(r *MonkeSockets.Room, c *MonkeSockets.Client, message []byte) {
	// 	r.Broadcast("chat:", message)
	// 	r.Broadcast("gun:", message)
	// 	r.Broadcast("test:", message)
	// }

	room.On("chat:", func(r *MonkeSockets.Room, c *MonkeSockets.Client, message []byte) {
		r.Broadcast("chat:", message)
		r.Broadcast("gun:", message)
		r.Broadcast("test:", message)
	})

	go room.Run()

	e.GET("/ws", room.WebSocket)
	//SOCKETS

	e.Logger.Fatal(e.Start(":3000"))
}
