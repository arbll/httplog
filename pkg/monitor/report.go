package monitor

import (
	"strings"

	"github.com/omen-/httplog/pkg/commonformat"
)

// TraficReport represents a report of the trafic for a given time frame.
type TraficReport struct {
	// RequestCount is the number of requests received
	RequestCount int64
	// BytesSentCount is the number of bytes sent from the server to the clients
	BytesSentCount int64
	// RequestsBySection contains the number of request for each sections
	RequestsBySection map[string]int64
}

func newTraficReport() *TraficReport {
	return &TraficReport{
		RequestsBySection: make(map[string]int64),
	}
}

// updateTraficReport updates a TraficReport with a new log entry.
func (traficReport *TraficReport) updateTraficReport(newLogEntry commonformat.LogEntry) {
	section := extractSection(newLogEntry.Request.Resource)
	traficReport.RequestsBySection[section]++
	traficReport.BytesSentCount += newLogEntry.BytesSent
	traficReport.RequestCount++
}

func extractSection(resource string) string {
	if strings.Count(resource, "/") < 2 {
		return resource
	}
	return resource[:strings.Index(resource[1:], "/")+1]
}
