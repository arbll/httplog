package commonformat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	var time, _ = time.Parse(TimeLayout, "28/Jul/2006:10:27:32 -0300")

	logEntry := LogEntry{
		IP:         "127.0.0.1",
		Identity:   "-",
		UserID:     "-",
		Time:       time,
		Request:    Request{Method: "GET", Resource: "/foo/bar", HTTPVersion: "HTTP/1.0"},
		StatusCode: 404,
		BytesSent:  int64(7218),
	}

	serializer := LogSerializer{}

	serializedEntry := serializer.SerializeEntry(logEntry)

	assert.Equal(t, `127.0.0.1 - - [28/Jul/2006:10:27:32 -0300] "GET /foo/bar HTTP/1.0" 404 7218`, serializedEntry)
}
