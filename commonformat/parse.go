package commonformat

import (
	"regexp"
	"strconv"

	"github.com/omen-/httplog"
)

const commonFormatPropertyCount = 7

var commonFormatRegexp = regexp.MustCompile(`^(\S+) (\S+) (\S+) \[(.*)\] "(.*)" (\d{3}) (\d+)`)

type LogParser struct{}

func (l LogParser) ParseLine(line string) (httplog.LogEntry, error) {

	var logEntry httplog.LogEntry
	catchedProperties := commonFormatRegexp.FindStringSubmatch(line)
	if len(catchedProperties) != commonFormatPropertyCount+1 {
		return logEntry, httplog.NewLogParseError(line)
	}

	logEntry.IP = catchedProperties[1]
	logEntry.Identity = catchedProperties[2]
	logEntry.UserID = catchedProperties[3]
	logEntry.DateTime = catchedProperties[4]
	logEntry.Request = catchedProperties[5]

	statusCode, err := strconv.Atoi(catchedProperties[6])
	if err != nil {
		return logEntry, httplog.NewLogParseError(line)
	}
	logEntry.StatusCode = statusCode

	bytesSent, err := strconv.ParseInt(catchedProperties[7], 10, 64)
	if err != nil {
		return logEntry, httplog.NewLogParseError(line)
	}
	logEntry.BytesSent = bytesSent

	return logEntry, nil
}
