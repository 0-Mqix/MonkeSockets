// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)

func home(e echo.Context) error {
	return e.HTML(200, "")
}

func (room *Room) ws(e echo.Context) error {
	serveWs(room, e)
	return nil
}

var events = make(map[string]func(*Room, []byte))

type SocketMessageJsonIn struct {
	Msg  string `json:"msg"`
	Test string `json:"test"`
}

func main() {
	e := echo.New()
	e.File("/", "home.html")

	room := newRoom()

	events["msg"] = func(r *Room, message []byte) {
		for client := range r.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(r.clients, client)
			}
		}
	}

	events["json"] = func(r *Room, message []byte) {
		split := bytes.SplitAfterN(message, []byte(":"), 2)
		event := string(split[0])
		if event == "chat"+":" {
			data := SocketMessageJsonIn{}
			json.Unmarshal(split[1], &data)
			fmt.Println(data.Msg, data.Test)
		}
	}

	go room.run(&events)

	e.GET("/ws", room.ws)

	e.Logger.Fatal(e.Start(":3000"))
}
