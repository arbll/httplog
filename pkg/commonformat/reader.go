package commonformat

import (
	"io"
	"log"

	"github.com/hpcloud/tail"
)

// Reader is in charge of reading lines from an actively written to common log
// format file.
type Reader struct {
	Logs       chan LogEntry
	parser     LogParser
	tailReader *tail.Tail
}

// NewReader returns a running common log format file reader.
func NewReader(filePath string, logParser LogParser) (*Reader, error) {
	var reader Reader

	locationEnd := &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}
	tailReader, err := tail.TailFile(filePath, tail.Config{Follow: true, ReOpen: true, Location: locationEnd, Logger: tail.DiscardingLogger})
	if err != nil {
		return &reader, err
	}

	reader.parser = logParser
	reader.tailReader = tailReader
	reader.Logs = make(chan LogEntry)

	go reader.handleRawLines(tailReader)

	return &reader, nil
}

// Close closes the reader.
func (r *Reader) Close() {
	r.tailReader.Stop()
	r.tailReader.Cleanup()
}

// handleRawLines reads and parse new entries.
func (r *Reader) handleRawLines(tailReader *tail.Tail) {
	for line := range tailReader.Lines {
		logEntry, err := r.parser.ParseLine(line.Text)
		if err != nil {
			log.Println(err)
			return
		}
		r.Logs <- logEntry
	}
}
