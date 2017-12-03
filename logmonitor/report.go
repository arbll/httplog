package logmonitor

import "github.com/omen-/httplog"
import "strings"

type TraficReport struct {
	RequestCount      int64
	BytesSentCount    int64
	RequestsBySection map[string]int64
}

func (traficReport *TraficReport) updateTraficReport(newLogEntry httplog.LogEntry) {
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
