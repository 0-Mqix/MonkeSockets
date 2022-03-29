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
			"onJoin:":  func(r *Room, c *Client, message []byte) {},
			"onLeave:": func(r *Room, c *Client, message []byte) {},
		},
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			r.Events["onJoin:"](r, client, nil)

		case client := <-r.unregister:
			r.Events["onLeave:"](r, client, nil)
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

func (r *Room) Broadcast(message []byte) {
	for c := range r.clients {
		c.SendMessage(message)
	}
}

func (r *Room) WebSocket(e echo.Context) error {
	ServeWs(r, e)
	return nil
}
