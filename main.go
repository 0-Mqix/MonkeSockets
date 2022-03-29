package main

import (
	"encoding/json"
	ms "projects/chat/MonkeSockets"

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

	//SOCKETS
	room := ms.NewRoom()

	room.Events["chat:"] = func(r *ms.Room, c *ms.Client, message []byte) {
		data := SocketMessageJsonIn{}
		json.Unmarshal(message, &data)
		r.Broadcast([]byte(data.Msg))
	}

	go room.Run()
	e.GET("/ws", room.WebSocket)
	//SOCKETS

	e.Logger.Fatal(e.Start(":3000"))
}
