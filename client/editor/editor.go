package editor

import (
	"sync"

	"github.com/nsf/termbox-go"
)

// EditorConfig holds any editor-specific configuration settings.
type EditorConfig struct {
	ScrollEnabled bool
}

// Editor represents the editor's skeleton.
// The editor is composed of two components:
// 1. an editable text area (primary interactive area).
// 2. a status bar (displays messages like user-join events).
type Editor struct {
	// Text contains the editor's content.
	Text []rune

	// Cursor represents the cursor position of the editor.
	Cursor int

	// Width represents the terminal's width in characters.
	Width int

	// Height represents the terminal's width in characters.
	Height int

	// ColOff is the number of columns between the start of a line and the left of the editor window.
	ColOff int

	// RowOff is the number of rows between the beginning of the text and the top of the editor window.
	RowOff int

	// ShowMsg acts like a switch for the status bar.
	ShowMsg bool

	// StatusMsg holds the text to be displayed in the status bar.
	StatusMsg string

	// StatusChan is used to send and receive status messages.
	StatusChan chan string

	// StatusMu protects against concurrent reads/writes to status bar info.
	StatusMu sync.Mutex

	// Users holds the names of all users connected to the server, displayed in the status bar.
	Users []string

	// ScrollEnabled determines whether or not the user can scroll past the initial
	// editor window. It is set by the EditorConfig.
	ScrollEnabled bool

	// IsConnected shows whether the editor is currently connected to the server.
	IsConnected bool

	// DrawChan is used to send and receive signals to update the terminal display.
	DrawChan chan int

	// mu prevents concurrent reads and writes to the editor state.
	mu sync.RWMutex
}

// userColors is a simple list of termbox attributes, used for coloring user names.
var userColors = []termbox.Attribute{
	termbox.ColorGreen,
	termbox.ColorYellow,
	termbox.ColorBlue,
	termbox.ColorMagenta,
	termbox.ColorCyan,
	termbox.ColorLightYellow,
	termbox.ColorLightMagenta,
	termbox.ColorLightGreen,
	termbox.ColorLightRed,
	termbox.ColorRed,
}

// NewEditor returns a new instance of the editor.
func NewEditor(conf EditorConfig) *Editor {
	return &Editor{
		ScrollEnabled: conf.ScrollEnabled,
		StatusChan:    make(chan string, 100),
		DrawChan:      make(chan int, 10000),
	}
}
