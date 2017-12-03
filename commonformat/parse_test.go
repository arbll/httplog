package commonformat

import (
	"fmt"
	"testing"
	"time"

	"github.com/omen-/httplog"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	var (
		validIP         = "127.0.0.1"
		validIdentity   = "-"
		validUserID     = "-"
		validTime, _    = time.Parse(TimeLayout, "28/Jul/2006:10:27:32 -0300")
		validRequest    = httplog.Request{Method: "GET", Resource: "/foo/bar", HTTPVersion: "HTTP/1.0"}
		validStatusCode = 404
		validBytesSent  = int64(7218)
	)

	validLogLine := fmt.Sprintf(`%v %v %v [%v] "%v" %v %v`, validIP, validIdentity, validUserID, validTime.Format(TimeLayout), validRequest.String(), validStatusCode, validBytesSent)

	var parser = LogParser{}

	logEntry, err := parser.ParseLine(validLogLine)
	if assert.Nil(t, err) {
		assert.Equal(t, validIP, logEntry.IP)
		assert.Equal(t, validIdentity, logEntry.Identity)
		assert.Equal(t, validUserID, logEntry.UserID)
		assert.Equal(t, validTime, logEntry.Time)
		assert.Equal(t, validRequest, logEntry.Request)
		assert.Equal(t, validStatusCode, logEntry.StatusCode)
		assert.Equal(t, validBytesSent, logEntry.BytesSent)
	}
}

func TestParseLineError(t *testing.T) {
	const invalidLogLine = "127.0.0.1 - -"

	var parser = LogParser{}

	_, err := parser.ParseLine(invalidLogLine)
	assert.NotNil(t, err)
}
