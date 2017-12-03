package commonformat

import (
	"fmt"
	"time"
)

const (
	TimeLayout = "02/Jan/2006:15:04:05 -0700"
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
