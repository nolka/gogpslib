package writer

import (
	"bytes"
	"log"
	"os"
)

func CreateFileWriter(fileName string) *FileWriter{
	w := &FileWriter{}
	w.FileName = fileName
	return w
}

type FileWriter struct {
	FileName string
	buf bytes.Buffer
}

func (s *FileWriter) Append(line string) {
	s.buf.WriteString(line)
}

func (s *FileWriter) Write() {
	f, err := os.Create(s.FileName)
	if err != nil {
		log.Print(err)
	}
	defer  f.Close()
	f.WriteString(s.buf.String())
}

