package commonformat

import (
	"os"
)

// Writer is in charge of writing lines to a common log format file.
type Writer struct {
	logSerializer *LogSerializer
	filePath      string
}

// NewWriter returns a common log format file writer.
func NewWriter(filePath string, logSerializer *LogSerializer) *Writer {
	return &Writer{
		filePath:      filePath,
		logSerializer: logSerializer,
	}
}

// WriteLogEntry writes an entry to the log file.
func (writer *Writer) WriteLogEntry(logEntry LogEntry) error {
	logFile, err := os.OpenFile(writer.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer logFile.Close()
	serializedLogEntry := writer.logSerializer.SerializeEntry(logEntry)
	_, err = logFile.WriteString(serializedLogEntry + "\n")
	return err
}
