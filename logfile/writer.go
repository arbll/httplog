package logfile

import (
	"os"

	"github.com/omen-/httplog"
)

type Writer struct {
	logSerializer httplog.LogSerializer
	filePath      string
}

func NewWriter(filePath string, logSerializer httplog.LogSerializer) Writer {
	return Writer{
		filePath:      filePath,
		logSerializer: logSerializer,
	}
}

//FIXME: Opening once is probably better if possible. Think to add close() and defer it
func (writer Writer) WriteLogEntry(logEntry httplog.LogEntry) error {
	logFile, err := os.OpenFile(writer.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	serializedLogEntry := writer.logSerializer.SerializeEntry(logEntry)
	_, err = logFile.WriteString(serializedLogEntry + "\n")
	return err
}
