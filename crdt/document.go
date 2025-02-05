package crdt

// Document is composed of characters.
type Document struct {
	Characters []Character
}

// New returns an initialized document.
func New() Document {
	return Document{Characters: []Character{CharacterStart, CharacterEnd}}
}

// Length returns the length of the document.
func (doc *Document) Length() int {
	return len(doc.Characters)
}

// ElementAt returns the character present in the position.
func (doc *Document) ElementAt(position int) (Character, error) {
	if position < 0 || position >= doc.Length() {
		return Character{}, ErrPositionOutOfBounds
	}
	return doc.Characters[position], nil
}

// Position returns the position (1-based) of the character with given ID.
// If not found, returns -1.
func (doc *Document) Position(charID string) int {
	for position, char := range doc.Characters {
		if charID == char.ID {
			return position + 1
		}
	}
	return -1
}

// Contains checks if a character (by ID) is present in the document.
func (doc *Document) Contains(charID string) bool {
	position := doc.Position(charID)
	return position != -1
}

// Find returns the character with the given ID, or a sentinel if not found.
func (doc *Document) Find(id string) Character {
	for _, char := range doc.Characters {
		if char.ID == id {
			return char
		}
	}
	return Character{ID: "-1"}
}

// Subseq returns the content between two positions (exclusive).
// wcharacterStart and wcharacterEnd should be valid markers in the Document.
func (doc *Document) Subseq(wcharacterStart, wcharacterEnd Character) ([]Character, error) {
	startPosition := doc.Position(wcharacterStart.ID)
	endPosition := doc.Position(wcharacterEnd.ID)

	if startPosition == -1 || endPosition == -1 {
		return doc.Characters, ErrBoundsNotPresent
	}
	if startPosition > endPosition {
		return doc.Characters, ErrBoundsNotPresent
	}
	if startPosition == endPosition {
		return []Character{}, nil
	}

	return doc.Characters[startPosition : endPosition-1], nil
}

// Left returns the ID of the character to the left of charID.
func (doc *Document) Left(charID string) string {
	i := doc.Position(charID)
	if i <= 0 {
		return doc.Characters[i].ID
	}
	return doc.Characters[i-1].ID
}

// Right returns the ID of the character to the right of charID.
func (doc *Document) Right(charID string) string {
	i := doc.Position(charID)
	if i >= len(doc.Characters)-1 {
		return doc.Characters[i-1].ID
	}
	return doc.Characters[i+1].ID
}

// SetText appends the characters from newDoc to the existing Document.
func (doc *Document) SetText(newDoc Document) {
	for _, char := range newDoc.Characters {
		c := Character{
			ID:         char.ID,
			Visible:    char.Visible,
			Value:      char.Value,
			IDPrevious: char.IDPrevious,
			IDNext:     char.IDNext,
		}
		doc.Characters = append(doc.Characters, c)
	}
}

// Content returns the visible content of the document as a string.
func Content(doc Document) string {
	value := ""
	for _, char := range doc.Characters {
		if char.Visible {
			value += char.Value
		}
	}
	return value
}

// IthVisible returns the ith visible character (1-based) in the document.
func IthVisible(doc Document, position int) Character {
	count := 0
	for _, char := range doc.Characters {
		if char.Visible {
			if count == position-1 {
				return char
			}
			count++
		}
	}
	return Character{ID: "-1"}
}
