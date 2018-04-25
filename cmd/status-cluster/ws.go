package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/layout"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	upgrader websocket.Upgrader
	hub      []*websocket.Conn

	Positions []*position
	layout    layout.Layout
	graph     *graph.Graph
}

func NewWSServer(layout layout.Layout) *WSServer {
	ws := &WSServer{
		upgrader: websocket.Upgrader{},
		layout:   layout,
	}
	ws.updatePositions()
	return ws
}

type WSResponse struct {
	Type      MsgType         `json:"type"`
	Positions []*position     `json:"positions,omitempty"`
	Graph     json.RawMessage `json:"graph,omitempty"`
}

type WSRequest struct {
	Cmd WSCommand `json:"cmd"`
}

type MsgType string
type WSCommand string

// WebSocket response types
const (
	RespPositions MsgType = "positions"
	RespGraph     MsgType = "graph"
)

// WebSocket commands
const (
	CmdInit    WSCommand = "init"
	CmdRefresh WSCommand = "refresh"
)

func (ws *WSServer) Handle(w http.ResponseWriter, r *http.Request) {
	c, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	ws.hub = append(ws.hub, c)

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", mt, err)
			break
		}
		ws.processRequest(c, mt, message)
	}
}

func (ws *WSServer) processRequest(c *websocket.Conn, mtype int, data []byte) {
	var cmd WSRequest
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		log.Fatal("unmarshal command", err)
		return
	}

	switch cmd.Cmd {
	case CmdInit:
		ws.sendGraphData(c)
		ws.updatePositions()
		ws.sendPositions(c)
	case CmdRefresh:
		ws.sendGraphData(c)
		ws.updatePositions()
		ws.sendPositions(c)
	}
}

func (ws *WSServer) sendMsg(c *websocket.Conn, msg *WSResponse) {
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
