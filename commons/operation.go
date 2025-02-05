package commons

// Operation defines a CRDT operation that can be applied to modify the document
type Operation struct {
	// Type indicates whether this is an insert or delete operation
	Type string `json:"type"`

	// Position specifies where in the document this operation occurs
	Position int `json:"position"`

	// Value contains the actual character or content being operated on
	Value string `json:"value"`
}
