package main

import "github.com/gorilla/websocket"

type position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	Z int32 `json:"z"`
}

func (ws *WSServer) sendPositions(c *websocket.Conn) {
	msg := &WSResponse{
		Type:      RespPositions,
		Idx:       ws.currentIdx,
		Positions: ws.history[ws.currentIdx].Positions,
	}

	ws.sendMsg(c, msg)
}

func (ws *WSServer) updateForcesAndPositions() {
	// positions
	nodes := ws.layout.Nodes()
	positions := []*position{}
	for i := 0; i < len(nodes); i++ {
		pos := &position{
			X: nodes[i].X,
			Y: nodes[i].Y,
			Z: nodes[i].Z,
		}
		positions = append(positions, pos)
	}

	// forces
	forces := ws.layout.Forces()

	hist := &ForceAndPosition{
		Positions: positions,
		Forces:    forces,
	}
	ws.history = append(ws.history, hist)
	ws.currentIdx = len(ws.history) - 1
	ws.broadcastPositions()
	ws.broadcastForces()
}

func (ws *WSServer) broadcastPositions() {
	for i := 0; i < len(ws.hub); i++ {
		ws.sendPositions(ws.hub[i])
	}
}
