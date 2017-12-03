package logmonitor

import "github.com/omen-/httplog"
import "strings"

type TraficReport struct {
	requestCount      int64
	bytesSentCount    int64
	requestsBySection map[string]int64
}

func (traficReport *TraficReport) updateTraficReport(newLogEntry httplog.LogEntry) {
	section := extractSection(newLogEntry.Request.Resource)
	traficReport.requestsBySection[section]++
	traficReport.bytesSentCount += newLogEntry.BytesSent
	traficReport.requestCount++
}

func extractSection(resource string) string {
	if strings.Count(resource, "/") < 2 {
		return resource
	}
	return resource[:strings.Index(resource[1:], "/")+1]
}
