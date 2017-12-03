package commonformat

import (
	"os"
)

type Writer struct {
	logSerializer LogSerializer
	filePath      string
}

func NewWriter(filePath string, logSerializer LogSerializer) Writer {
	return Writer{
		filePath:      filePath,
		logSerializer: logSerializer,
	}
}

//FIXME: Opening once is probably better if possible. Think to add close() and defer it
func (writer Writer) WriteLogEntry(logEntry LogEntry) error {
	logFile, err := os.OpenFile(writer.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	serializedLogEntry := writer.logSerializer.SerializeEntry(logEntry)
	_, err = logFile.WriteString(serializedLogEntry + "\n")
	return err
}
