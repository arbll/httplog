package monitor

import (
	"strings"

	"github.com/omen-/httplog/pkg/commonformat"
)

// TrafficReport represents a report of the traffic for a given time frame.
type TrafficReport struct {
	// RequestCount is the number of requests received
	RequestCount int64
	// BytesSentCount is the number of bytes sent from the server to the clients
	BytesSentCount int64
	// RequestsBySection contains the number of requests for each sections
	RequestsBySection map[string]int64
}

func newTrafficReport() *TrafficReport {
	return &TrafficReport{
		RequestsBySection: make(map[string]int64),
	}
}

// updateTrafficReport updates a TrafficReport with a new log entry.
func (trafficReport *TrafficReport) updateTrafficReport(newLogEntry commonformat.LogEntry) {
	section := extractSection(newLogEntry.Request.Resource)
	trafficReport.RequestsBySection[section]++
	trafficReport.BytesSentCount += newLogEntry.BytesSent
	trafficReport.RequestCount++
}

func extractSection(resource string) string {
	if strings.Count(resource, "/") < 2 {
		return resource
	}
	return resource[:strings.Index(resource[1:], "/")+1]
}
