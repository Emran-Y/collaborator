package editor

// GetText returns the editor's content.
func (e *Editor) GetText() []rune {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Text
}

// SetText sets the given string as the editor's content.
func (e *Editor) SetText(text string) {
	e.mu.Lock()
	e.Text = []rune(text)
	e.mu.Unlock()
}

// GetX returns the X-axis component of the current cursor position.
func (e *Editor) GetX() int {
	x, _ := e.calcXY(e.Cursor)
	return x
}

// SetX sets the X-axis component of the current cursor position to the specified X position.
func (e *Editor) SetX(x int) {
	e.Cursor = x
}

// GetY returns the Y-axis component of the current cursor position.
func (e *Editor) GetY() int {
	_, y := e.calcXY(e.Cursor)
	return y
}

// GetWidth returns the editor's width (in characters).
func (e *Editor) GetWidth() int {
	return e.Width
}

// GetHeight returns the editor's height (in characters).
func (e *Editor) GetHeight() int {
	return e.Height
}

// SetSize sets the editor size to the specific width and height.
func (e *Editor) SetSize(w, h int) {
	e.Width = w
	e.Height = h
}

// GetRowOff returns the vertical offset of the editor window from the start of the text.
func (e *Editor) GetRowOff() int {
	return e.RowOff
}

// GetColOff returns the horizontal offset of the editor window from the start of a line.
func (e *Editor) GetColOff() int {
	return e.ColOff
}

// IncRowOff increments the vertical offset of the editor window from the start of the text by inc.
func (e *Editor) IncRowOff(inc int) {
	e.RowOff += inc
}

// IncColOff increments the horizontal offset of the editor window from the start of a line by inc.
func (e *Editor) IncColOff(inc int) {
	e.ColOff += inc
}

// SendDraw sends a draw signal to the drawLoop. Use this function to ensure concurrency safety for rendering.
func (e *Editor) SendDraw() {
	e.DrawChan <- 1
}
