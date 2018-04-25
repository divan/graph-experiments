package main

import "github.com/gorilla/websocket"

func (ws *WSServer) sendStats(c *websocket.Conn) {
	msg := &WSResponse{
		Type:  RespStats,
		Stats: makeStats(ws.graph),
	}

	ws.sendMsg(c, msg)
}

func (ws *WSServer) broadcastStats() {
	for i := 0; i < len(ws.hub); i++ {
		ws.sendStats(ws.hub[i])
	}
}
