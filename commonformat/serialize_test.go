package commonformat

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/omen-/httplog"
)

func TestSerialize(t *testing.T) {
	logEntry := httplog.LogEntry{
		IP:         "127.0.0.1",
		Identity:   "-",
		UserID:     "-",
		DateTime:   "28/Jul/2006:10:27:32 -0300",
		Request:    "GET /foo/bar HTTP/1.0",
		StatusCode: 404,
		BytesSent:  int64(7218),
	}

	serializer := LogSerializer{}

	serializedEntry := serializer.SerializeEntry(logEntry)

	assert.Equal(t, `127.0.0.1 - - [28/Jul/2006:10:27:32 -0300] "GET /foo/bar HTTP/1.0" 404 7218`, serializedEntry)
}
