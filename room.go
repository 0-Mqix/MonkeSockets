// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
//Modified by MqixSchool

package MonkeSockets

import (
	"bytes"

	"github.com/labstack/echo/v4"
)

type Room struct {
	clients    map[*Client]bool
	message    chan SocketMessage
	register   chan *Client
	unregister chan *Client
	Events     map[string]func(*Room, *Client, []byte)
}

func NewRoom() *Room {
	return &Room{
		message:    make(chan SocketMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),

		//map were you can inject your own functions in that can use the room, client, and the message for all the incoming messages
		Events: map[string]func(*Room, *Client, []byte){
			"join:":  func(r *Room, c *Client, message []byte) {},
			"leave:": func(r *Room, c *Client, message []byte) {},
		},
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			r.Events["join:"](r, client, nil)

		case client := <-r.unregister:
			r.Events["leave:"](r, client, nil)
			_, ok := r.clients[client]

			if ok {
				delete(r.clients, client)
				close(client.send)
			}

		case s := <-r.message:
			split := bytes.SplitAfterN(s.Message, []byte(":"), 2)
			event := string(split[0])
			function, ok := r.Events[event]
			if ok {
				function(r, s.Client, split[1])
			}
		}
	}
}

func (r *Room) On(event string, function func(*Room, *Client, []byte)) {
	r.Events[event] = function
}

func (r *Room) Broadcast(event string, message []byte) {
	for c := range r.clients {
		c.SendMessage(event, message)
	}
}

// serveWs handles websocket requests from the peer.
func (room *Room) WebSocket(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return nil
	}
	client := &Client{room: room, conn: conn, send: make(chan []byte, 256)}
	client.room.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	return nil
}
