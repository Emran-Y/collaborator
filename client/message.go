package main

import (
	"fmt"
	"strconv"
	"strings"

	"terminal_collab/commons"
	"terminal_collab/crdt"

	"github.com/gorilla/websocket"
)

// handleMsg updates the CRDT document with the contents of the message.
func handleMsg(msg commons.Message, conn *websocket.Conn) {
	switch msg.Type {
	case commons.DocSyncMessage:
		logger.Infof("DOCSYNC RECEIVED, updating local doc %+v\n", msg.Document)
		doc = msg.Document
		e.SetText(crdt.Content(doc))

	case commons.DocReqMessage:
		logger.Infof("DOCREQ RECEIVED, sending local document to %v\n", msg.ID)
		docMsg := commons.Message{Type: commons.DocSyncMessage, Document: doc, ID: msg.ID}
		_ = conn.WriteJSON(&docMsg)

	case commons.SiteIDMessage:
		siteID, err := strconv.Atoi(msg.Text)
		if err != nil {
			logger.Errorf("failed to set siteID, err: %v\n", err)
		}
		crdt.SiteID = siteID
		logger.Infof("SITE ID %v, INTENDED SITE ID: %v", crdt.SiteID, siteID)

	case commons.JoinMessage:
		e.StatusChan <- fmt.Sprintf("%s has joined the session!", msg.Username)

	case commons.UsersMessage:
		e.StatusMu.Lock()
		e.Users = strings.Split(msg.Text, ",")
		e.StatusMu.Unlock()

	default:
		switch msg.Operation.Type {
		case "insert":
			_, err := doc.Insert(msg.Operation.Position, msg.Operation.Value)
			if err != nil {
				logger.Errorf("failed to insert, err: %v\n", err)
			}

			e.SetText(crdt.Content(doc))
			if msg.Operation.Position-1 <= e.Cursor {
				e.MoveCursor(len(msg.Operation.Value), 0)
			}
			logger.Infof("REMOTE INSERT: %s at position %v\n", msg.Operation.Value, msg.Operation.Position)

		case "delete":
			_ = doc.Delete(msg.Operation.Position)
			e.SetText(crdt.Content(doc))
			if msg.Operation.Position-1 <= e.Cursor {
				e.MoveCursor(-len(msg.Operation.Value), 0)
			}
			logger.Infof("REMOTE DELETE: position %v\n", msg.Operation.Position)
		}
	}

	// Debug function (you mentioned it's toggled via a flag).
	printDoc(doc)

	e.SendDraw()
}

// getMsgChan returns a message channel that repeatedly reads from a websocket connection.
func getMsgChan(conn *websocket.Conn) chan commons.Message {
	messageChan := make(chan commons.Message)
	go func() {
		for {
			var msg commons.Message

			// Read message.
			err := conn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Errorf("websocket error: %v", err)
				}
				e.IsConnected = false
				e.StatusChan <- "lost connection!"
				break
			}

			logger.Infof("message received: %+v\n", msg)

			// send message through channel
			messageChan <- msg
		}
	}()
	return messageChan
}
