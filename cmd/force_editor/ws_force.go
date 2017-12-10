package main

import "github.com/gorilla/websocket"

type force struct {
	FX float64 `json:"fx"`
	FY float64 `json:"fy"`
	FZ float64 `json:"fz"`
}

func (ws *WSServer) broadcastForces() {
	for i := 0; i < len(ws.hub); i++ {
		ws.sendForces(ws.hub[i])
	}
}

func (ws *WSServer) sendForces(c *websocket.Conn) {
	msg := &WSResponse{
		Type:   RespForces,
		Idx:    ws.currentIdx,
		Forces: ws.history[ws.currentIdx].Forces,
	}

	ws.sendMsg(c, msg)
}
