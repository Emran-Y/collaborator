package main

import (
	"time"
)

// handleStatusMsg asynchronously waits for messages from e.StatusChan and
// displays the message when it arrives.
func handleStatusMsg() {
	for msg := range e.StatusChan {
		e.StatusMu.Lock()
		e.StatusMsg = msg
		e.ShowMsg = true
		e.StatusMu.Unlock()

		logger.Infof("got status message: %s", e.StatusMsg)

		e.SendDraw()
		time.Sleep(3 * time.Second)

		e.StatusMu.Lock()
		e.ShowMsg = false
		e.StatusMu.Unlock()

		e.SendDraw()
	}
}
