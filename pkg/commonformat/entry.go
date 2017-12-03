package commonformat

import (
	"fmt"
	"time"
)

const (
	// TimeLayout is the common log format time layout.
	TimeLayout = "02/Jan/2006:15:04:05 -0700"
)

// LogEntry represents a common log format entry.
type LogEntry struct {
	IP         string
	Identity   string
	UserID     string
	Time       time.Time
	Request    Request
	StatusCode int
	BytesSent  int64
}

// Request represents an http request in common log format.
type Request struct {
	Method      string
	Resource    string
	HTTPVersion string
}

func (r *Request) String() string {
	return fmt.Sprintf("%v %v %v", r.Method, r.Resource, r.HTTPVersion)
}
