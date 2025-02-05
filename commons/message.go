package commons

import (
	"terminal_collab/crdt"

	"github.com/google/uuid"
)

// Message defines the structure for messages exchanged between clients and server
type Message struct {
	Username string `json:"username"`

	// Text contains the message body used for join notifications, siteID assignment,
	// and active user lists
	Text string `json:"text"`

	// Type indicates what kind of message this is (see MessageType constants)
	Type MessageType `json:"type"`

	// ID uniquely identifies the client
	ID uuid.UUID `json:"ID"`

	// Operation contains any CRDT operations to be applied
	Operation Operation `json:"operation"`

	// Document holds the full CRDT document state
	// Note: Only sent when absolutely necessary due to size
	Document crdt.Document `json:"document"`
}

// MessageType categorizes the different kinds of messages that can be sent
type MessageType string

// Supported message types:
// DocSyncMessage - Synchronize document state between clients
// DocReqMessage - Request latest document state
// SiteIDMessage - Assign unique site ID to client
// JoinMessage - Notify when clients join
// UsersMessage - Update list of active users

const (
	DocSyncMessage MessageType = "docSync"
	DocReqMessage  MessageType = "docReq"
	SiteIDMessage  MessageType = "SiteID"
	JoinMessage    MessageType = "join"
	UsersMessage   MessageType = "users"
)
