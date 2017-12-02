package commonformat

import (
	"fmt"

	"github.com/omen-/httplog"
)

type LogSerializer struct{}

func (l LogSerializer) SerializeEntry(logEntry httplog.LogEntry) string {
	return fmt.Sprintf(`%v %v %v [%v] "%v" %v %v`, logEntry.IP, logEntry.Identity,
		logEntry.UserID, logEntry.DateTime, logEntry.Request, logEntry.StatusCode, logEntry.BytesSent)
}
