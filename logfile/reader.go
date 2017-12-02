package logfile

import (
	"log"

	"github.com/hpcloud/tail"
	"github.com/omen-/httplog"
)

type Reader struct {
	parser httplog.LogParser
	logs   chan httplog.LogEntry
}

func NewReader(filePath string, logParser httplog.LogParser) (Reader, error) {
	var reader Reader

	tailReader, err := tail.TailFile(filePath, tail.Config{Follow: true})
	if err != nil {
		return reader, err
	}

	reader.parser = logParser
	reader.logs = make(chan httplog.LogEntry)

	go reader.handleRawLines(tailReader) //FIXME : Context / Close

	return reader, nil
}

func (r Reader) Logs() chan httplog.LogEntry {
	return r.logs
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
