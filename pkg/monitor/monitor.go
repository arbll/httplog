// Package monitor allows common log format file monitoring with periodic
// traffic reports and alerts.
package monitor

import (
	"log"
	"time"

	"github.com/omen-/httplog/pkg/commonformat"
)

const (
	trafficReportBufferSize = 4096
	trafficAlertBufferSize  = 4096
)

// Config holds configuration options that can be passed to New in order to
// customize the monitor.
type Config struct {
	// TrafficReportPeriod is the duration between each report.
	TrafficReportPeriod time.Duration
	// AlertPeriod is the duration the monitor will keep logs to check if an Alert should be raised.
	AlertPeriod time.Duration
	// AlertThreshold is the threshold after which an allert will be triggered.
	AlertThreshold int64
}

// Monitor is in charge of monitoring a common log format file.
type Monitor struct {
	logReader      *commonformat.Reader
	trafficReports chan TrafficReport
	trafficAlerts  chan Alert
	config         Config
}

// New returns a new monitor for the given log file that will generate alerts
// at the given threshold.
func New(config Config, logReader *commonformat.Reader) *Monitor {
	monitor := Monitor{
		config:         config,
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
// By default, the channel has a buffer size of 4096.
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

// monitorLogs is the event loop used to monitor the log file.
func (monitor *Monitor) monitorLogs() {
	logsChannel := monitor.logReader.Logs

	alertMonitor := newAlertMonitor(monitor.config.AlertPeriod, monitor.config.AlertThreshold)

	trafficReportTicker := time.NewTicker(monitor.config.TrafficReportPeriod)
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
