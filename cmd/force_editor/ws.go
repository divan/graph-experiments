package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Msg struct {
	Type MsgType `json:"type"`
	Data []byte  `json:"data"`
}

type MsgType string

const Positions MsgType = "positions"

func handleWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", mt, err)
			break
		}
		log.Printf("recv: %s", message)

		sendPositions(c)
	}
}

func sendPositions(c *websocket.Conn) {
	msg := &Msg{
		Type: Positions,
		Data: []byte("test"),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}

	err = c.WriteMessage(1, data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}
