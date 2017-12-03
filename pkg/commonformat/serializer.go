package commonformat

import (
	"fmt"
)

// LogSerializer is in charge of serializing common log format entries.
type LogSerializer struct{}

// SerializeEntry serializes a common log format entry.
func (l *LogSerializer) SerializeEntry(logEntry LogEntry) string {
	return fmt.Sprintf(`%v %v %v [%v] "%v" %v %v`, logEntry.IP, logEntry.Identity,
		logEntry.UserID, logEntry.Time.Format(TimeLayout), logEntry.Request.String(), logEntry.StatusCode, logEntry.BytesSent)
}
