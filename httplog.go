package httplog

type LogEntry struct {
	IP         string
	Identity   string
	UserID     string
	DateTime   string
	Request    string
	StatusCode int
	BytesSent  int64
}

type LogParser interface {
	parseLine(line string) (LogEntry, error)
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
	serializeEntry(logEntry LogEntry) (string, error)
}
