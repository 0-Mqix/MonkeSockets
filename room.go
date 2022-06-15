// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
//Modified by MqixSchool

package MonkeSockets

import (
	"bytes"
)

type Room struct {
	clients    map[*Client]bool
	message    chan SocketMessage
	Register   chan *Client
	unregister chan *Client
	Events     map[string]func(*Room, *Client, []byte)
}

func New() *Room {
	return &Room{
		Register:   make(chan *Client),
		message:    make(chan SocketMessage),
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
		case client := <-r.Register:
			r.clients[client] = true
			r.Events["join:"](r, client, nil)

		case client := <-r.unregister:
			r.Events["leave:"](r, client, nil)
			delete(r.clients, client)

		case s := <-r.message:
			split := bytes.SplitAfterN(s.Message, []byte(":"), 2)
			function, ok := r.Events[string(split[0])]
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
		c.Send(event, message)
	}
}
