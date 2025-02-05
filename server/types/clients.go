package types

import (
	"sync"

	"github.com/google/uuid"
)

// Clients stores and manages all connected clients.
type Clients struct {
	// List of active clients, keyed by UUID.
	List map[uuid.UUID]*Client

	// Protects reads/writes to the 'List'.
	Mu sync.RWMutex

	// Channels for concurrency control (add, read, delete, name updates).
	DeleteRequests     chan DeleteRequest
	ReadRequests       chan ReadRequest
	AddRequests        chan *Client
	NameUpdateRequests chan NameUpdate
}

// NewClients returns a new instance of a Clients struct.
func NewClients() *Clients {
	return &Clients{
		List:               make(map[uuid.UUID]*Client),
		Mu:                 sync.RWMutex{},
		DeleteRequests:     make(chan DeleteRequest),
		ReadRequests:       make(chan ReadRequest, 10000),
		AddRequests:        make(chan *Client),
		NameUpdateRequests: make(chan NameUpdate),
	}
}

// Client holds the information of a connected user.
type Client struct {
	// Actual WebSocket connection
	Conn   interface{} // Will set it to *websocket.Conn at runtime
	SiteID string
	ID     uuid.UUID

	// Protects against concurrent writes to a single WebSocket Conn
	WriteMu sync.Mutex

	// Protects data races on a Client's internal fields
	Mu sync.Mutex

	Username string
}

// DeleteRequest indicates that a client should be removed from 'List'.
type DeleteRequest struct {
	ID   uuid.UUID
	Done chan int
}

// ReadRequest is used to retrieve either all clients or a single client.
type ReadRequest struct {
	ReadAll bool
	ID      uuid.UUID
	Resp    chan *Client
}

// NameUpdate is used to update a Client's username.
type NameUpdate struct {
	ID      uuid.UUID
	NewName string
}
