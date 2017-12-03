package commonformat

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

const commonFormatPropertyCount = 7
const requestPropertyCount = 3

var commonFormatRegexp = regexp.MustCompile(`^(\S+) (\S+) (\S+) \[(.*)\] "(.*)" (\d{3}) (\d+)`)
var requestRegexp = regexp.MustCompile(`^(\w+) (\S+) (\S+)`)

type LogParser struct{}

type ParseError struct {
	line string
}

func (e *ParseError) Error() string {
	return e.line
}

func (l LogParser) ParseLine(line string) (LogEntry, error) {
	var logEntry LogEntry
	matchedGroups := commonFormatRegexp.FindStringSubmatch(line)
	if len(matchedGroups) != commonFormatPropertyCount+1 {
		return logEntry, &ParseError{line}
	}

	logEntry.IP = matchedGroups[1]
	logEntry.Identity = matchedGroups[2]
	logEntry.UserID = matchedGroups[3]

	time, err := time.Parse(TimeLayout, matchedGroups[4])
	if err != nil {
		return logEntry, &ParseError{line}
	}
	logEntry.Time = time

	request, err := parseRequest(matchedGroups[5])
	if err != nil {
		return logEntry, &ParseError{line}
	}
	logEntry.Request = request

	statusCode, err := strconv.Atoi(matchedGroups[6])
	if err != nil {
		return logEntry, &ParseError{line}
	}
	logEntry.StatusCode = statusCode

	bytesSent, err := strconv.ParseInt(matchedGroups[7], 10, 64)
	if err != nil {
		return logEntry, &ParseError{line}
	}
	logEntry.BytesSent = bytesSent

	return logEntry, nil
}

func parseRequest(rawRequest string) (Request, error) {
	var request Request
	matchedGroups := requestRegexp.FindStringSubmatch(rawRequest)
	if len(matchedGroups) != requestPropertyCount+1 {
		return request, errors.New("")
	}

	request.Method = matchedGroups[1]
	request.Resource = matchedGroups[2]
	request.HTTPVersion = matchedGroups[3]

	return request, nil
}
