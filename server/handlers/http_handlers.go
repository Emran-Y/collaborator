package handlers

import (
	"net/http"
	"strconv"
	"time"

	"terminal_collab/commons"
	"terminal_collab/server/helper_function"
	"terminal_collab/server/types"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// HandleConn upgrades an HTTP request to WebSocket, sets up the new client, and reads messages.
func HandleConn(w http.ResponseWriter, r *http.Request, c *types.Clients) {
	conn, err := types.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		color.Red("Error upgrading connection to websocket: %v\n", err)
		return
	}
	defer conn.Close()

	clientID := uuid.New()

	// Protect increments to SiteID with a mutex
	types.Mu.Lock()
	types.SiteID++
	newSiteID := types.SiteID
	types.Mu.Unlock()

	// Create new client
	cl := &types.Client{
		Conn:     conn, // store the actual *websocket.Conn here
		SiteID:   strconv.Itoa(newSiteID),
		ID:       clientID,
		Username: "",
	}

	// Register the client with our global Clients struct
	helper_function.Add(c, cl)

	// Send that client its own siteID
	siteIDMsg := commons.Message{
		Type: commons.SiteIDMessage,
		Text: cl.SiteID,
		ID:   clientID,
	}
	if err := helper_function.BroadcastOne(c, siteIDMsg, clientID); err != nil {
		color.Red("Error broadcasting site ID to client: %v", err)
		return
	}

	// Ask another client for the doc
	docReq := commons.Message{
		Type: commons.DocReqMessage,
		ID:   clientID,
	}
	helper_function.BroadcastOneExcept(c, docReq, clientID)

	// Send the list of usernames to all
	helper_function.SendUsernames(c)

	// Loop to read messages from the client connection
	for {
		var msg commons.Message
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				color.Red("Unexpected close from %s: %v", cl.Username, err)
			}
			color.Red("Client disconnected: %v", cl.Username)
			helper_function.Delete(c, cl.ID)
			return
		}

		// If it's a DocSyncMessage, pass it to the sync channel (before rewriting msg.ID).
		if msg.Type == commons.DocSyncMessage {
			types.SyncChan <- msg
			continue
		}

		// Reassign the message ID to the sender's ID
		msg.ID = clientID

		// Send message to messageChan for logging/broadcast
		types.MessageChan <- msg
	}
}

// HandleMsg listens on MessageChan and broadcasts incoming messages to other clients.
func HandleMsg(c *types.Clients) {
	for {
		msg := <-types.MessageChan
		t := time.Now().Format(time.ANSIC)

		switch msg.Type {
		case commons.JoinMessage:
			helper_function.UpdateName(c, msg.ID, msg.Username)
			color.Green("%s >> %s %s (ID: %s)\n", t, msg.Username, msg.Text, msg.ID)
			helper_function.SendUsernames(c)

		case "operation":
			color.Green("Operation >> %+v from ID=%s\n", msg.Operation, msg.ID)

		default:
			color.Green("%s >> Unknown message type: %v\n", t, msg)
			helper_function.SendUsernames(c)
			continue
		}
		// Broadcast to all except the sender
		helper_function.BroadcastAllExcept(c, msg, msg.ID)
	}
}

// HandleSync reads from the SyncChan and routes doc sync messages to the correct client(s).
func HandleSync(c *types.Clients) {
	for {
		syncMsg := <-types.SyncChan
		switch syncMsg.Type {
		case commons.DocSyncMessage:
			// Broadcast this doc sync to the indicated client ID
			helper_function.BroadcastOne(c, syncMsg, syncMsg.ID)

		case commons.UsersMessage:
			color.Blue("Usernames: %s", syncMsg.Text)
			helper_function.BroadcastAll(c, syncMsg)
		}
	}
}
