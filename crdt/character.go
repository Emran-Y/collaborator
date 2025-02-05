package crdt

// Character represents a character in the document.
// As per section 3.1, Data Model in the paper (https://hal.inria.fr/inria-00108523/document)
type Character struct {
	ID         string
	Visible    bool
	Value      string
	IDPrevious string
	IDNext     string
}
