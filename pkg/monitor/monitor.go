// Package monitor allows common log format file monitoring with periodic
// traffic reports and alerts.
package monitor

import (
	"log"
	"time"

	"github.com/omen-/httplog/pkg/commonformat"
)

const (
	trafficReportPeriod = 10 * time.Second
	alertMonitorPeriod  = 2 * time.Minute

	trafficReportBufferSize = 4096
	trafficAlertBufferSize  = 4096
)

// Monitor is in charge of monitoring a common log format file.
type Monitor struct {
	treshold       int64
	logReader      *commonformat.Reader
	trafficReports chan TrafficReport
	trafficAlerts  chan Alert
}

// New returns a new monitor for the given log file that will generate alerts
// at the given treshold.
func New(treshold int64, logReader *commonformat.Reader) *Monitor {
	monitor := Monitor{
		treshold:       treshold,
		logReader:      logReader,
		trafficReports: make(chan TrafficReport, trafficReportBufferSize),
		trafficAlerts:  make(chan Alert, trafficAlertBufferSize),
	}

	go monitor.monitorLogs()

	return &monitor
}

// TrafficReports return a channel to the traffic reports.
// By default, the channel has a buffer size of 4096.
func (monitor *Monitor) TrafficReports() chan TrafficReport {
	return monitor.trafficReports
}

// Alerts return a channel to the traffic alerts.
// By default, the channel has a buffer size of 4096
func (monitor *Monitor) Alerts() chan Alert {
	return monitor.trafficAlerts
}

func (monitor *Monitor) onTrafficReport(trafficReport *TrafficReport) {
	select {
	case monitor.trafficReports <- *trafficReport:
	default:
		log.Println("Traffic report buffer is full, discarding new report")
	}
}

func (monitor *Monitor) onTrafficAlert(alert Alert) {
	select {
	case monitor.trafficAlerts <- alert:
	default:
		log.Println("Traffic alerts buffer is full, discarding new alert")
	}
}

// monitorLogs is the event loop used to monitor the log file
func (monitor *Monitor) monitorLogs() {
	logsChannel := monitor.logReader.Logs

	alertMonitor := newAlertMonitor(alertMonitorPeriod, monitor.treshold)

	trafficReportTicker := time.NewTicker(trafficReportPeriod)
	trafficReport := newTrafficReport()

	for {
		select {
		case logEntry := <-logsChannel:
			trafficReport.updateTrafficReport(logEntry)
			alert := alertMonitor.addLogEntry(logEntry)
			if alert != nil {
				monitor.onTrafficAlert(alert)
			}
		case <-trafficReportTicker.C:
			monitor.onTrafficReport(trafficReport)
			trafficReport = newTrafficReport()
		}
	}
}
