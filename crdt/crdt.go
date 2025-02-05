package crdt

import "fmt"

type CRDT interface {
	Insert(position int, value string) (string, error)
	Delete(position int) string
}

func IsCRDT(c CRDT) {
	// Test function to verify CRDT functionality
	fmt.Println(c.Insert(1, "a"))
}
