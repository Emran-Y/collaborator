package crdt

import (
	"errors"
	"sync"
)

// mu protects global CRDT-related state.
var mu sync.Mutex

// SiteID is a globally unique variable used with the local clock to generate
// identifiers for characters in the document.
var SiteID = 0

// LocalClock is incremented whenever an insert operation takes place.
// It is used to uniquely identify each character.
var LocalClock = 0

// CharacterStart is placed at the start of the document.
var CharacterStart = Character{ID: "start", Visible: false, Value: "", IDPrevious: "", IDNext: "end"}

// CharacterEnd is placed at the end of the document.
var CharacterEnd = Character{ID: "end", Visible: false, Value: "", IDPrevious: "start", IDNext: ""}

// Errors used by the CRDT.
var (
	ErrPositionOutOfBounds = errors.New("position out of bounds")
	ErrEmptyWCharacter     = errors.New("empty char ID provided")
	ErrBoundsNotPresent    = errors.New("subsequence bound(s) not present")
)
