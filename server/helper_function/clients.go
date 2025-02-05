package helper_function

import (
	"errors"

	"terminal_collab/commons"
	"terminal_collab/server/types"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// HandleClients runs in a goroutine and processes add/read/delete/name-update requests.
// This replaces the old (c *Clients) handle() method, turning it into a free function.
func HandleClients(c *types.Clients) {
	for {
		select {
		// Delete request
		case req := <-c.DeleteRequests:
			closeClient(c, req.ID)
			req.Done <- 1
			close(req.Done)

		// Read request (all clients or single client)
		case req := <-c.ReadRequests:
			if req.ReadAll {
				c.Mu.RLock()
				for _, client := range c.List {
					req.Resp <- client
				}
				c.Mu.RUnlock()
				close(req.Resp)
			} else {
				c.Mu.RLock()
				cl := c.List[req.ID]
				c.Mu.RUnlock()
				req.Resp <- cl
				close(req.Resp)
			}

		// Add request
		case cl := <-c.AddRequests:
			c.Mu.Lock()
			c.List[cl.ID] = cl
			c.Mu.Unlock()

		// Name update request
		case n := <-c.NameUpdateRequests:
			c.Mu.RLock()
			cl, found := c.List[n.ID]
			c.Mu.RUnlock()
			if found {
				cl.Mu.Lock()
				cl.Username = n.NewName
				cl.Mu.Unlock()
			}
		}
	}
}

// Add a new client
func Add(c *types.Clients, client *types.Client) {
	c.AddRequests <- client
}

// Delete a client
func Delete(c *types.Clients, id uuid.UUID) {
	req := types.DeleteRequest{
		ID:   id,
		Done: make(chan int),
	}
	c.DeleteRequests <- req
	<-req.Done
	SendUsernames(c)
}

// UpdateName updates a client’s username
func UpdateName(c *types.Clients, id uuid.UUID, newName string) {
	c.NameUpdateRequests <- types.NameUpdate{
		ID:      id,
		NewName: newName,
	}
}

// GetAll returns a channel of *Client references (all clients).
func GetAll(c *types.Clients) chan *types.Client {
	c.Mu.RLock()
	resp := make(chan *types.Client, len(c.List))
	c.Mu.RUnlock()

	c.ReadRequests <- types.ReadRequest{ReadAll: true, Resp: resp}
	return resp
}

// Get returns a channel with the single requested client.
func Get(c *types.Clients, id uuid.UUID) chan *types.Client {
	resp := make(chan *types.Client)
	c.ReadRequests <- types.ReadRequest{ReadAll: false, ID: id, Resp: resp}
	return resp
}

// BroadcastAll sends 'msg' to all active clients.
func BroadcastAll(c *types.Clients, msg commons.Message) {
	color.Blue("sending message to all users. Text: %s", msg.Text)
	for cl := range GetAll(c) {
		if err := send(cl, msg); err != nil {
			color.Red("Error broadcasting to client %v: %v", cl.ID, err)
			Delete(c, cl.ID)
		}
	}
}

// BroadcastAllExcept sends 'msg' to all but the given UUID.
func BroadcastAllExcept(c *types.Clients, msg commons.Message, except uuid.UUID) {
	for cl := range GetAll(c) {
		if cl.ID == except {
			continue
		}
		if err := send(cl, msg); err != nil {
			color.Red("Error broadcasting to client %v: %v", cl.ID, err)
			Delete(c, cl.ID)
		}
	}
}

// BroadcastOne sends 'msg' to exactly one client with ID=dst.
func BroadcastOne(c *types.Clients, msg commons.Message, dst uuid.UUID) error {
	cl := <-Get(c, dst)
	if cl == nil {
		return errors.New("no client found with given ID")
	}
	if err := send(cl, msg); err != nil {
		color.Red("Error sending to client %v: %v", cl.ID, err)
		Delete(c, cl.ID)
		return err
	}
	return nil
}

// BroadcastOneExcept sends a message to exactly one client who is NOT 'except'.
func BroadcastOneExcept(c *types.Clients, msg commons.Message, except uuid.UUID) {
	for cl := range GetAll(c) {
		if cl.ID == except {
			continue
		}
		if err := send(cl, msg); err != nil {
			color.Red("Error sending to client %v: %v", cl.ID, err)
			Delete(c, cl.ID)
			continue
		}
		break // Stop after sending to first available client
	}
}

// SendUsernames builds a comma-separated list of client usernames and sends it to SyncChan.
func SendUsernames(c *types.Clients) {
	var users string
	for cl := range GetAll(c) {
		users += cl.Username + ","
	}
	types.SyncChan <- commons.Message{Text: users, Type: commons.UsersMessage}
}

// -------------------------
// Internal helper functions
// -------------------------

// closeClient closes the underlying websocket Conn for a given client ID.
func closeClient(c *types.Clients, id uuid.UUID) {
	c.Mu.RLock()
	cl, ok := c.List[id]
	c.Mu.RUnlock()

	if ok && cl.Conn != nil {
		ws, ok2 := cl.Conn.(*websocket.Conn)
		if ok2 {
			if err := ws.Close(); err != nil {
				color.Red("Error closing websocket for client %v: %v", cl.ID, err)
			}
		}
		color.Red("Removing client: %v (%s)", cl.ID, cl.Username)
	} else {
		color.Red("Could not find client %v to close.", id)
		return
	}

	c.Mu.Lock()
	delete(c.List, id)
	c.Mu.Unlock()
}

// send writes JSON to a client’s websocket, guarded by WriteMu.
func send(cl *types.Client, data interface{}) error {
	cl.WriteMu.Lock()
	defer cl.WriteMu.Unlock()

	ws, ok := cl.Conn.(*websocket.Conn)
	if !ok {
		return errors.New("client.Conn is not a websocket.Conn")
	}
	return ws.WriteJSON(data)
}
