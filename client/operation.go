package main

import (
	"terminal_collab/commons"

	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
)

// Constants for the CRDT operation type.
const (
	OperationInsert = iota
	OperationDelete
)

// performOperation performs a CRDT insert or delete operation on the local document
// and sends a message over the WebSocket connection.
func performOperation(opType int, ev termbox.Event, conn *websocket.Conn) {
	// Get position and value.
	ch := string(ev.Ch)

	var msg commons.Message

	// Modify local state (CRDT) first.
	switch opType {
	case OperationInsert:
		logger.Infof("LOCAL INSERT: %s at cursor position %v\n", ch, e.Cursor)

		text, err := doc.Insert(e.Cursor+1, ch)
		if err != nil {
			e.SetText(text)
			logger.Errorf("CRDT error: %v\n", err)
		}
		e.SetText(text)

		e.MoveCursor(1, 0)
		msg = commons.Message{
			Type: "operation",
			Operation: commons.Operation{
				Type:     "insert",
				Position: e.Cursor,
				Value:    ch,
			},
		}

	case OperationDelete:
		logger.Infof("LOCAL DELETE: cursor position %v\n", e.Cursor)

		if e.Cursor-1 < 0 {
			e.Cursor = 0
		}

		text := doc.Delete(e.Cursor)
		e.SetText(text)

		msg = commons.Message{
			Type: "operation",
			Operation: commons.Operation{
				Type:     "delete",
				Position: e.Cursor,
			},
		}
		e.MoveCursor(-1, 0)
	}

	// Send the message.
	if e.IsConnected {
		err := conn.WriteJSON(msg)
		if err != nil {
			e.IsConnected = false
			e.StatusChan <- "lost connection!"
		}
	}
}
