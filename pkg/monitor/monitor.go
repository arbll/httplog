// Package monitor allows common log format file monitoring with periodic
// trafic reports and alerts.
package monitor

import (
	"log"
	"time"

	"github.com/omen-/httplog/pkg/commonformat"
)

const (
	traficReportPeriod = 10 * time.Second
	alertMonitorPeriod = 2 * time.Minute

	traficReportBufferSize = 4096
	traficAlertBufferSize  = 4096
)

// Monitor is in charge of monitoring a common log format file.
type Monitor struct {
	treshold      int64
	logReader     *commonformat.Reader
	traficReports chan TraficReport
	traficAlerts  chan Alert
}

// New returns a new monitor for the given log file that will generate alerts
// at the given treshold.
func New(treshold int64, logReader *commonformat.Reader) *Monitor {
	monitor := Monitor{
		treshold:      treshold,
		logReader:     logReader,
		traficReports: make(chan TraficReport, traficReportBufferSize),
		traficAlerts:  make(chan Alert, traficAlertBufferSize),
	}

	go monitor.monitorLogs()

	return &monitor
}

// TraficReports return a channel to the trafic reports.
// By default, the channel has a buffer size of 4096.
func (monitor *Monitor) TraficReports() chan TraficReport {
	return monitor.traficReports
}

// Alerts return a channel to the trafic alerts.
// By default, the channel has a buffer size of 4096
func (monitor *Monitor) Alerts() chan Alert {
	return monitor.traficAlerts
}

func (monitor *Monitor) onTraficReport(traficReport *TraficReport) {
	select {
	case monitor.traficReports <- *traficReport:
	default:
		log.Println("Trafic report buffer is full, discarding new report")
	}
}

func (monitor *Monitor) onTraficAlert(alert Alert) {
	select {
	case monitor.traficAlerts <- alert:
	default:
		log.Println("Trafic alerts buffer is full, discarding new alert")
	}
}

// monitorLogs is the event loop used to monitor the log file
func (monitor *Monitor) monitorLogs() {
	logsChannel := monitor.logReader.Logs

	alertMonitor := newAlertMonitor(alertMonitorPeriod, monitor.treshold)

	traficReportTicker := time.NewTicker(traficReportPeriod)
	traficReport := newTraficReport()

	for {
		select {
		case logEntry := <-logsChannel:
			traficReport.updateTraficReport(logEntry)
			alert := alertMonitor.addLogEntry(logEntry)
			if alert != nil {
				monitor.onTraficAlert(alert)
			}
		case <-traficReportTicker.C:
			monitor.onTraficReport(traficReport)
			traficReport = newTraficReport()
		}
	}
}
