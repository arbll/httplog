package httplog

type LogParser interface {
	ParseLine(line string) (LogEntry, error)
}

type LogParseError struct {
	parsedLine string
}

func (e *LogParseError) Error() string {
	return e.parsedLine
}

func NewLogParseError(parsedLine string) error {
	return &LogParseError{parsedLine}
}

type LogSerializer interface {
	SerializeEntry(logEntry LogEntry) string
}
