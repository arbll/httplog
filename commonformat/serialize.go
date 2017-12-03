package commonformat

import (
	"fmt"
)

type LogSerializer struct{}

func (l LogSerializer) SerializeEntry(logEntry LogEntry) string {
	return fmt.Sprintf(`%v %v %v [%v] "%v" %v %v`, logEntry.IP, logEntry.Identity,
		logEntry.UserID, logEntry.Time.Format(TimeLayout), logEntry.Request.String(), logEntry.StatusCode, logEntry.BytesSent)
}
