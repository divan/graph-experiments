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

	layout layout.LayoutWithDebug
	graph  *graph.Data

	history    []*ForceAndPosition
	currentIdx int
}

type ForceAndPosition struct {
	Positions []*position
	Forces    layout.ForcesDebugData
}

func NewWSServer(layout layout.LayoutWithDebug) *WSServer {
	ws := &WSServer{
		upgrader: websocket.Upgrader{},
		layout:   layout,
	}
	ws.updateForcesAndPositions()
	return ws
}

type WSResponse struct {
	Type      MsgType                `json:"type"`
	Idx       int                    `json:"idx"`
	Positions []*position            `json:"positions,omitempty"`
	Graph     *graph.Data            `json:"graph,omitempty"`
	Forces    layout.ForcesDebugData `json:"forces,omitempty"`
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
	RespForces    MsgType = "forces"
)

// WebSocket commands
const (
	CmdInit  WSCommand = "init"
	CmdPrev            = "prev"
	CmdNext            = "next"
	CmdCalc            = "calc"
	CmdReset           = "reset"
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
		ws.sendPositions(c)
	case CmdPrev:
		ws.cmdPrev()
	case CmdNext:
		ws.cmdNext()
	case CmdCalc:
		ws.cmdCalc()
	case CmdReset:
		ws.cmdReset()
	}
}

func (ws *WSServer) cmdPrev() {
	if ws.currentIdx > 0 {
		ws.currentIdx--
		ws.broadcastPositions()
		ws.broadcastForces()
	}
}

func (ws *WSServer) cmdNext() {
	if ws.currentIdx == len(ws.history)-1 {
		ws.layout.CalculateN(1)
		ws.updateForcesAndPositions()
	} else {
		ws.currentIdx++
		ws.broadcastPositions()
		ws.broadcastForces()
	}
}

func (ws *WSServer) cmdCalc() {
	ws.layout.Calculate()
	ws.updateForcesAndPositions()
}

func (ws *WSServer) cmdReset() {
	ws.layout.Reset()
	ws.updateForcesAndPositions()
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
