package types

import (
	"sync"

	"terminal_collab/commons"

	"github.com/gorilla/websocket"
)

// Here are your "constants" (though some are actually mutable vars).
// They live in the 'types' package so other packages can import them.

var (
	// Monotonically increasing site ID, unique to each client.
	SiteID = 0

	// Mutex for protecting site ID increments.
	Mu sync.Mutex

	// Upgrader instance to upgrade all HTTP connections to a WebSocket.
	Upgrader = websocket.Upgrader{}

	// Channel for client messages.
	MessageChan = make(chan commons.Message)

	// Channel for document sync messages.
	SyncChan = make(chan commons.Message)
)
