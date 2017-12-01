package commonformat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	const validIP = "127.0.0.1"
	const validIdentity = "-"
	const validUserID = "-"
	const validDateTime = "28/Jul/2006:10:27:32 -0300"
	const validRequest = "GET /foo/bar HTTP/1.0"
	const validStatusCode = 404
	const validBytesSent = int64(7218)

	validLogLine := fmt.Sprintf(`%v %v %v [%v] "%v" %v %v`, validIP, validIdentity, validUserID, validDateTime, validRequest, validStatusCode, validBytesSent)

	var parser = LogParser{}

	logEntry, err := parser.ParseLine(validLogLine)
	if assert.Nil(t, err) {
		assert.Equal(t, validIP, logEntry.IP)
		assert.Equal(t, validIdentity, logEntry.Identity)
		assert.Equal(t, validUserID, logEntry.UserID)
		assert.Equal(t, validDateTime, logEntry.DateTime)
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
