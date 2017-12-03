package httplog

import (
	"fmt"
	"time"
)

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
