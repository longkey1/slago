package output

import (
	"os"
)

// StdoutWriter writes to stdout
type StdoutWriter struct {
	*JSONWriter
}

// NewStdoutWriter creates a new stdout writer
func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{
		JSONWriter: NewJSONWriter(os.Stdout, true),
	}
}
