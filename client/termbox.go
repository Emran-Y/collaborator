package main

import (
	"errors"
	"fmt"

	"terminal_collab/commons"
	"terminal_collab/crdt"

	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
	"github.com/sirupsen/logrus"
)

// handleTermboxEvent handles key input by updating the local CRDT document
// and sending a message over the WebSocket connection.
func handleTermboxEvent(ev termbox.Event, conn *websocket.Conn) error {
	if ev.Type == termbox.EventKey {
		switch ev.Key {
		case termbox.KeyEsc, termbox.KeyCtrlC:
			return errors.New("pairpad: exiting")

		case termbox.KeyCtrlS:
			if fileName == "" {
				fileName = "pairpad-content.txt"
			}

			err := crdt.Save(fileName, &doc)
			if err != nil {
				logrus.Errorf("Failed to save to %s", fileName)
				e.StatusChan <- fmt.Sprintf("Failed to save to %s", fileName)
				return err
			}

			e.StatusChan <- fmt.Sprintf("Saved document to %s", fileName)

		case termbox.KeyCtrlL:
			if fileName != "" {
				logger.Log(logrus.InfoLevel, "LOADING DOCUMENT")
				newDoc, err := crdt.Load(fileName)
				if err != nil {
					logrus.Errorf("failed to load file %s", fileName)
					e.StatusChan <- fmt.Sprintf("Failed to load %s", fileName)
					return err
				}
				e.StatusChan <- fmt.Sprintf("Loading %s", fileName)
				doc = newDoc
				e.SetX(0)
				e.SetText(crdt.Content(doc))

				logger.Log(logrus.InfoLevel, "SENDING DOCUMENT")
				docMsg := commons.Message{Type: commons.DocSyncMessage, Document: doc}
				_ = conn.WriteJSON(&docMsg)
			} else {
				e.StatusChan <- "No file to load!"
			}

		case termbox.KeyArrowLeft, termbox.KeyCtrlB:
			e.MoveCursor(-1, 0)

		case termbox.KeyArrowRight, termbox.KeyCtrlF:
			e.MoveCursor(1, 0)

		case termbox.KeyArrowUp, termbox.KeyCtrlP:
			e.MoveCursor(0, -1)

		case termbox.KeyArrowDown, termbox.KeyCtrlN:
			e.MoveCursor(0, 1)

		case termbox.KeyHome:
			e.SetX(0)

		case termbox.KeyEnd:
			e.SetX(len(e.Text))

		case termbox.KeyBackspace, termbox.KeyBackspace2:
			performOperation(OperationDelete, ev, conn)
		case termbox.KeyDelete:
			performOperation(OperationDelete, ev, conn)

		case termbox.KeyTab:
			for i := 0; i < 4; i++ {
				ev.Ch = ' '
				performOperation(OperationInsert, ev, conn)
			}

		case termbox.KeyEnter:
			ev.Ch = '\n'
			performOperation(OperationInsert, ev, conn)

		case termbox.KeySpace:
			ev.Ch = ' '
			performOperation(OperationInsert, ev, conn)

		default:
			if ev.Ch != 0 {
				performOperation(OperationInsert, ev, conn)
			}
		}
	}

	e.SendDraw()
	return nil
}

// getTermboxChan returns a channel of termbox Events repeatedly waiting on user input.
func getTermboxChan() chan termbox.Event {
	termboxChan := make(chan termbox.Event)
	go func() {
		for {
			termboxChan <- termbox.PollEvent()
		}
	}()
	return termboxChan
}
