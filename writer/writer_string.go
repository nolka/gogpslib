package writer

import (
	"bytes"
)

func CreateStringWriter() *StringWriter {
	return &StringWriter{}
}

type Writer interface {
	Append(line string)
	Write()
}

type StringWriter struct {
	Content string
	buf bytes.Buffer
}

func (s *StringWriter) Append(line string) {
	 s.buf.WriteString(line)
}

func (s *StringWriter) Write() {
	s.Content = s.buf.String()
}
