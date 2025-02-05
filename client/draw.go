package main

// drawLoop continuously waits for draw signals, then re-renders the editor.
func drawLoop() {
	for {
		<-e.DrawChan
		e.Draw()
	}
}
