package main

import "github.com/gorilla/websocket"

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
