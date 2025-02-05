package crdt

import (
	"fmt"
)

// LocalInsert inserts the character into the document at the given position.
func (doc *Document) LocalInsert(char Character, position int) (*Document, error) {
	if position <= 0 || position >= doc.Length() {
		return doc, ErrPositionOutOfBounds
	}
	if char.ID == "" {
		return doc, ErrEmptyWCharacter
	}

	doc.Characters = append(doc.Characters[:position],
		append([]Character{char}, doc.Characters[position:]...)...,
	)

	// Update next/previous pointers for neighbors.
	doc.Characters[position-1].IDNext = char.ID
	doc.Characters[position+1].IDPrevious = char.ID
	return doc, nil
}

// IntegrateInsert inserts 'char' based on 'charPrev' and 'charNext'.
func (doc *Document) IntegrateInsert(char, charPrev, charNext Character) (*Document, error) {
	subsequence, err := doc.Subseq(charPrev, charNext)
	if err != nil {
		return doc, err
	}

	// Position of the 'charNext' in doc
	position := doc.Position(charNext.ID)
	position--

	// If subsequence is empty, insert at current position.
	if len(subsequence) == 0 {
		return doc.LocalInsert(char, position)
	}

	// If subsequence has 1 char, insert at previous pos
	if len(subsequence) == 1 {
		return doc.LocalInsert(char, position-1)
	}

	// Otherwise, we find where to place 'char' by comparing IDs recursively.
	i := 1
	for i < len(subsequence)-1 && subsequence[i].ID < char.ID {
		i++
	}
	return doc.IntegrateInsert(char, subsequence[i-1], subsequence[i])
}

// GenerateInsert creates a new visible character for the given value
// and inserts it into the doc at the 1-based 'position'.
func (doc *Document) GenerateInsert(position int, value string) (*Document, error) {
	mu.Lock()
	LocalClock++
	mu.Unlock()

	charPrev := IthVisible(*doc, position-1)
	charNext := IthVisible(*doc, position)

	// Use defaults if not found
	if charPrev.ID == "-1" {
		charPrev = doc.Find("start")
	}
	if charNext.ID == "-1" {
		charNext = doc.Find("end")
	}

	char := Character{
		ID:         fmt.Sprint(SiteID) + fmt.Sprint(LocalClock),
		Visible:    true,
		Value:      value,
		IDPrevious: charPrev.ID,
		IDNext:     charNext.ID,
	}

	return doc.IntegrateInsert(char, charPrev, charNext)
}

// IntegrateDelete marks a character (by ID) invisible.
func (doc *Document) IntegrateDelete(char Character) *Document {
	position := doc.Position(char.ID)
	if position == -1 {
		return doc
	}
	doc.Characters[position-1].Visible = false
	return doc
}

// GenerateDelete finds the ith visible character and marks it invisible.
func (doc *Document) GenerateDelete(position int) *Document {
	char := IthVisible(*doc, position)
	return doc.IntegrateDelete(char)
}

// Insert is the CRDT interface to insert 'value' at 'position'.
func (doc *Document) Insert(position int, value string) (string, error) {
	newDoc, err := doc.GenerateInsert(position, value)
	if err != nil {
		return Content(*doc), err
	}
	return Content(*newDoc), nil
}

// Delete is the CRDT interface to delete the character at 'position'.
func (doc *Document) Delete(position int) string {
	newDoc := doc.GenerateDelete(position)
	return Content(*newDoc)
}
