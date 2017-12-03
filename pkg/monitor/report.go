package monitor

import (
	"strings"

	"github.com/omen-/httplog/pkg/commonformat"
)

type TraficReport struct {
	RequestCount      int64
	BytesSentCount    int64
	RequestsBySection map[string]int64
}

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
