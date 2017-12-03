package httplog

import "time"
import "fmt"

//TODO: Split everything in different files

type LogEntry struct {
	IP         string
	Identity   string
	UserID     string
	Time       time.Time
	Request    Request
	StatusCode int
	BytesSent  int64
}

type Request struct {
	Method      string
	Resource    string
	HTTPVersion string
}

func (r Request) String() string {
	return fmt.Sprintf("%v %v %v", r.Method, r.Resource, r.HTTPVersion)
}

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
