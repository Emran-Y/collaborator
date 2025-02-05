package crdt

import (
	"os"
	"strings"
)

// Load reads a text file from disk and converts it into a CRDT document.
func Load(fileName string) (Document, error) {
	doc := New()
	content, err := os.ReadFile(fileName)
	if err != nil {
		return doc, err
	}
	lines := strings.Split(string(content), "\n")
	pos := 1
	for i := 0; i < len(lines); i++ {
		for j := 0; j < len(lines[i]); j++ {
			_, err := doc.Insert(pos, string(lines[i][j]))
			if err != nil {
				return doc, err
			}
			pos++
		}
		if i < len(lines)-1 { // avoid insertion of '\n' on last line
			_, err := doc.Insert(pos, "\n")
			if err != nil {
				return doc, err
			}
			pos++
		}
	}
	return doc, nil
}

// Save writes data to the named file, creating it if necessary.
// The contents of the file are overwritten.
func Save(fileName string, doc *Document) error {
	return os.WriteFile(fileName, []byte(Content(*doc)), 0644)
}
