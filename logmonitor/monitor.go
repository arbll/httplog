package logmonitor

import (
	"log"
	"time"

	"github.com/omen-/httplog"
)

const (
	traficReportPeriod = 10 * time.Second
	alertMonitorPeriod = 1 * time.Minute

	traficReportBufferSize = 4096
	traficAlertBufferSize  = 4096
)

type Monitor struct {
	treshold      int64
	logStream     httplog.LogStream
	traficReports chan TraficReport
	traficAlerts  chan httplog.Alert
}

func New(treshold int64, logStream httplog.LogStream) *Monitor {
	monitor := Monitor{
		treshold:      treshold,
		logStream:     logStream,
		traficReports: make(chan TraficReport, traficReportBufferSize),
		traficAlerts:  make(chan httplog.Alert, traficAlertBufferSize),
	}

	go monitor.monitorLogs()

	return &monitor
}

func (monitor *Monitor) TraficReports() chan TraficReport {
	return monitor.traficReports
}

func (monitor *Monitor) Alerts() chan httplog.Alert {
	return monitor.traficAlerts
}

func (monitor *Monitor) onTraficReport(traficReport TraficReport) {
	select {
	case monitor.traficReports <- traficReport:
	default:
		log.Println("Trafic report buffer is full, discarding new report")
	}
}

func (monitor *Monitor) onTraficAlert(alert httplog.Alert) {
	select {
	case monitor.traficAlerts <- alert:
	default:
		log.Println("Trafic alerts buffer is full, discarding new alert")
	}
}

func (monitor *Monitor) monitorLogs() {
	logsChannel := monitor.logStream.Logs()

	alertMonitor := newAlertMonitor(alertMonitorPeriod, monitor.treshold)

	traficReportTicker := time.NewTicker(traficReportPeriod)
	traficReport := TraficReport{RequestsBySection: make(map[string]int64)}

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
			traficReport = TraficReport{RequestsBySection: make(map[string]int64)}
		}
	}
}
