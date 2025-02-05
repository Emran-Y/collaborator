package editor

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// Draw updates the UI by setting cells with the editor's content.
func (e *Editor) Draw() {
	_ = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	e.mu.RLock()
	cursor := e.Cursor
	e.mu.RUnlock()

	cx, cy := e.calcXY(cursor)

	// draw cursor x position relative to row offset
	if cx-e.GetColOff() > 0 {
		cx -= e.GetColOff()
	}

	// draw cursor y position relative to row offset
	if cy-e.GetRowOff() > 0 {
		cy -= e.GetRowOff()
	}

	termbox.SetCursor(cx-1, cy-1)

	// find the starting and ending row of the termbox window.
	yStart := e.GetRowOff()
	yEnd := yStart + e.GetHeight() - 1 // -1 accounts for the status bar

	// find the starting column of the termbox window.
	xStart := e.GetColOff()

	x, y := 0, 0
	e.mu.RLock()
	defer e.mu.RUnlock()

	for i := 0; i < len(e.Text) && y < yEnd; i++ {
		if e.Text[i] == rune('\n') {
			x = 0
			y++
		} else {
			// Set cell content. setX and setY account for the window offset.
			setY := y - yStart
			setX := x - xStart
			termbox.SetCell(setX, setY, e.Text[i], termbox.ColorDefault, termbox.ColorDefault)

			// Update x by rune's width.
			x = x + runewidth.RuneWidth(e.Text[i])
		}
	}

	e.DrawStatusBar()

	// Flush back buffer!
	termbox.Flush()
}

// DrawStatusBar shows all status and debug information on the bottom line of the editor.
func (e *Editor) DrawStatusBar() {
	e.StatusMu.Lock()
	showMsg := e.ShowMsg
	e.StatusMu.Unlock()
	if showMsg {
		e.DrawStatusMsg()
	} else {
		e.DrawInfoBar()
	}

	// Render connection indicator
	if e.IsConnected {
		termbox.SetBg(e.Width-1, e.Height-1, termbox.ColorGreen)
	} else {
		termbox.SetBg(e.Width-1, e.Height-1, termbox.ColorRed)
	}
}

// DrawStatusMsg draws the editor's status message at the bottom of the termbox window.
func (e *Editor) DrawStatusMsg() {
	e.StatusMu.Lock()
	statusMsg := e.StatusMsg
	e.StatusMu.Unlock()
	for i, r := range []rune(statusMsg) {
		termbox.SetCell(i, e.Height-1, r, termbox.ColorDefault, termbox.ColorDefault)
	}
}

// DrawInfoBar draws the editor's debug information and the names of the
// active users in the editing session at the bottom of the termbox window.
func (e *Editor) DrawInfoBar() {
	e.StatusMu.Lock()
	users := e.Users
	e.StatusMu.Unlock()

	length := len(e.Text)

	x := 0
	for i, user := range users {
		for _, r := range user {
			colorIdx := i % len(userColors)
			termbox.SetCell(x, e.Height-1, r, userColors[colorIdx], termbox.ColorDefault)
			x++
		}
		termbox.SetCell(x, e.Height-1, ' ', termbox.ColorDefault, termbox.ColorDefault)
		x++
	}

	e.mu.RLock()
	cursor := e.Cursor
	e.mu.RUnlock()

	cx, cy := e.calcXY(cursor)
	debugInfo := fmt.Sprintf(" x=%d, y=%d, cursor=%d, len(text)=%d", cx, cy, e.Cursor, length)

	for _, r := range debugInfo {
		termbox.SetCell(x, e.Height-1, r, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
}
