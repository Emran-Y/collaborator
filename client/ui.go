// Package main provides the main functionality for the collaborative terminal editor.
package main

import (
	"terminal_collab/client/editor"
	"terminal_collab/crdt"

	"github.com/gorilla/websocket"
	"github.com/nsf/termbox-go"
)

// UIConfig holds configuration options for the terminal user interface.
type UIConfig struct {
	EditorConfig editor.EditorConfig
}

// initUI initializes the terminal UI and editor with the given WebSocket connection
// and configuration. It sets up the editor state, starts background goroutines for
// status messages and screen drawing, and runs the main event loop.
func initUI(conn *websocket.Conn, conf UIConfig) error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()

	// Initialize and configure the editor
	e = editor.NewEditor(conf.EditorConfig)
	e.SetSize(termbox.Size())
	e.SetText(crdt.Content(doc))
	e.SendDraw()
	e.IsConnected = true

	// Start background goroutines
	go handleStatusMsg()
	go drawLoop()

	// Run the main event loop
	err = mainLoop(conn)
	if err != nil {
		return err
	}

	return nil
}

// mainLoop handles the main event loop of the editor, processing both terminal events
// from user input and incoming WebSocket messages from the server. It coordinates
// between local changes and remote updates to maintain consistency.
func mainLoop(conn *websocket.Conn) error {
	termboxChan := getTermboxChan()
	msgChan := getMsgChan(conn)

	for {
		select {
		case termboxEvent := <-termboxChan:
			err := handleTermboxEvent(termboxEvent, conn)
			if err != nil {
				return err
			}
		case msg := <-msgChan:
			handleMsg(msg, conn)
		}
	}
}
