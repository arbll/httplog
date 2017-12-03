package logfile

import (
	"io"
	"log"

	"github.com/hpcloud/tail"
	"github.com/omen-/httplog"
)

type Reader struct {
	parser     httplog.LogParser
	logs       chan httplog.LogEntry
	tailReader *tail.Tail
}

func NewReader(filePath string, logParser httplog.LogParser) (Reader, error) {
	var reader Reader

	locationEnd := &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}
	tailReader, err := tail.TailFile(filePath, tail.Config{Follow: true, ReOpen: true, Location: locationEnd, Logger: tail.DiscardingLogger})
	if err != nil {
		return reader, err
	}

	reader.parser = logParser
	reader.logs = make(chan httplog.LogEntry)
	reader.tailReader = tailReader

	go reader.handleRawLines(tailReader)

	return reader, nil
}

func (r Reader) Logs() chan httplog.LogEntry {
	return r.logs
}

func (r Reader) Close() {
	r.tailReader.Stop()
	r.tailReader.Cleanup()
}

func (r *Reader) handleRawLines(tailReader *tail.Tail) {
	for line := range tailReader.Lines {
		logEntry, err := r.parser.ParseLine(line.Text)
		if err != nil {
			log.Println(err)
			return
		}
		r.logs <- logEntry
	}
}
