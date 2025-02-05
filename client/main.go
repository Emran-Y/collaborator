package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"terminal_collab/client/editor"

	"terminal_collab/commons"
	"terminal_collab/crdt"

	"github.com/Pallinder/go-randomdata"
	"github.com/sirupsen/logrus"
)

// Global variables used throughout the application
var (
	// doc represents the collaborative document
	doc = crdt.New()

	// logger handles application logging
	logger = logrus.New()

	// e is the main editor instance
	e = editor.NewEditor(editor.EditorConfig{})

	// fileName stores the name of the file being edited
	fileName string

	// flags stores command line configuration
	flags Flags
)

func main() {
	// Parse command line flags
	flags = parseFlags()

	// Setup input scanner
	s := bufio.NewScanner(os.Stdin)

	// Generate random name or get from user input
	name := randomdata.SillyName()

	if flags.Login {
		fmt.Print("Enter your name: ")
		s.Scan()
		name = s.Text()
	}

	// Establish WebSocket connection
	conn, _, err := createConn(flags)
	if err != nil {
		fmt.Printf("Connection error, exiting: %s\n", err)
		return
	}
	defer conn.Close()

	// Send join message to server
	msg := commons.Message{Username: name, Text: "has joined the session.", Type: commons.JoinMessage}
	_ = conn.WriteJSON(msg)

	// Setup logging
	logFile, debugLogFile, err := setupLogger(logger)
	if err != nil {
		fmt.Printf("Failed to setup logger, exiting: %s\n", err)
		return
	}
	defer closeLogFiles(logFile, debugLogFile)

	// Load document from file if specified
	if flags.File != "" {
		if doc, err = crdt.Load(flags.File); err != nil {
			fmt.Printf("failed to load document: %s\n", err)
			return
		}
	}

	// Configure and initialize UI
	uiConfig := UIConfig{
		EditorConfig: editor.EditorConfig{
			ScrollEnabled: flags.Scroll,
		},
	}

	err = initUI(conn, uiConfig)
	if err != nil {
		if strings.HasPrefix(err.Error(), "pairpad") {
			fmt.Println("exiting session.")
			return
		}

		fmt.Printf("TUI error, exiting: %s\n", err)
		return
	}
}
