package editor

import "github.com/mattn/go-runewidth"

// MoveCursor updates the cursor position horizontally by a given x increment, and
// vertically by one line in the direction indicated by y. The positive directions are
// right and down, respectively.
// This is used by the UI layer, where it updates the cursor position on keypresses.
func (e *Editor) MoveCursor(x, y int) {
	if len(e.Text) == 0 && e.Cursor == 0 {
		return
	}
	// Move cursor horizontally.
	newCursor := e.Cursor + x

	// Move cursor vertically.
	if y > 0 {
		newCursor = e.calcCursorDown()
	}
	if y < 0 {
		newCursor = e.calcCursorUp()
	}

	if e.ScrollEnabled {
		cx, cy := e.calcXY(newCursor)

		// move the window to adjust for the cursor
		rowStart := e.GetRowOff()
		rowEnd := e.GetRowOff() + e.GetHeight() - 1

		if cy <= rowStart { // scroll up
			e.IncRowOff(cy - rowStart - 1)
		}
		if cy > rowEnd { // scroll down
			e.IncRowOff(cy - rowEnd)
		}

		colStart := e.GetColOff()
		colEnd := e.GetColOff() + e.GetWidth()

		if cx <= colStart { // scroll left
			e.IncColOff(cx - (colStart + 1))
		}
		if cx > colEnd { // scroll right
			e.IncColOff(cx - colEnd)
		}
	}

	// Reset to bounds.
	if newCursor > len(e.Text) {
		newCursor = len(e.Text)
	}
	if newCursor < 0 {
		newCursor = 0
	}

	e.mu.Lock()
	e.Cursor = newCursor
	e.mu.Unlock()
}

// For the functions calcCursorUp and calcCursorDown, newline characters are found by iterating
// backward and forward from the current cursor position. These characters are taken as the "start"
// and "end" of the current line. The "offset" from the start of the current line to the cursor is
// used for the final position on the target line, in case the target line is shorter or longer.

// calcCursorUp calculates and returns the intended cursor position after moving the cursor up one line.
func (e *Editor) calcCursorUp() int {
	pos := e.Cursor
	offset := 0

	// If the initial cursor is out of bounds or on a newline, move it.
	if pos == len(e.Text) || (pos >= 0 && e.Text[pos] == '\n') {
		offset++
		pos--
	}
	if pos < 0 {
		pos = 0
	}

	start, end := pos, pos

	// Find the start of the current line.
	for start > 0 && e.Text[start] != '\n' {
		start--
	}

	// If the cursor is on the first line, move to the beginning of the text.
	if start == 0 {
		return 0
	}

	// Find the end of the current line.
	for end < len(e.Text) && e.Text[end] != '\n' {
		end++
	}

	// Find the start of the previous line.
	prevStart := start - 1
	for prevStart >= 0 && e.Text[prevStart] != '\n' {
		prevStart--
	}

	// Calculate the distance from the start of this line to the cursor.
	offset += pos - start
	if offset <= start-prevStart {
		return prevStart + offset
	}
	return start
}

// calcCursorDown calculates and returns the intended cursor position after moving the cursor down one line.
func (e *Editor) calcCursorDown() int {
	pos := e.Cursor
	offset := 0

	// If the initial cursor position is out of bounds or on a newline, move it.
	if pos == len(e.Text) || (pos >= 0 && e.Text[pos] == '\n') {
		offset++
		pos--
	}
	if pos < 0 {
		pos = 0
	}

	start, end := pos, pos

	// Find the start of the current line.
	for start > 0 && e.Text[start] != '\n' {
		start--
	}

	// Special check for the first line if no newline at index 0.
	if start == 0 && e.Text[start] != '\n' {
		offset++
	}

	// Find the end of the current line.
	for end < len(e.Text) && e.Text[end] != '\n' {
		end++
	}

	// If on a newline, increment end so start != end.
	if pos < len(e.Text) && e.Text[pos] == '\n' && e.Cursor != 0 {
		end++
	}

	// If the Cursor is already on the last line, move to the end of the text.
	if end == len(e.Text) {
		return len(e.Text)
	}

	// Find the end of the next line.
	nextEnd := end + 1
	for nextEnd < len(e.Text) && e.Text[nextEnd] != '\n' {
		nextEnd++
	}

	// Calculate the distance from the start of the current line to the cursor.
	offset += pos - start
	if offset < nextEnd-end {
		return end + offset
	}
	return nextEnd
}

// calcXY returns the x and y coordinates of the cell at the given index in the text.
func (e *Editor) calcXY(index int) (int, int) {
	x := 1
	y := 1

	if index < 0 {
		return x, y
	}

	e.mu.RLock()
	length := len(e.Text)
	e.mu.RUnlock()

	if index > length {
		index = length
	}

	for i := 0; i < index; i++ {
		e.mu.RLock()
		r := e.Text[i]
		e.mu.RUnlock()

		if r == '\n' {
			x = 1
			y++
		} else {
			x = x + runewidth.RuneWidth(r)
		}
	}
	return x, y
}
