package logmonitor

import (
	"container/list"
	"time"

	"github.com/omen-/httplog"
)

type alertMonitor struct {
	period            time.Duration
	logList           *list.List
	totalTrafic       int64
	threshold         int64
	wasAboveThreshold bool
}

type TraficAlert struct {
	triggeredAt   time.Time
	alert         string
	averageTrafic int64
	direction     bool
}

func (traficAlert TraficAlert) TriggeredAt() time.Time {
	return traficAlert.triggeredAt
}

func (traficAlert TraficAlert) Alert() string {
	return traficAlert.alert
}

func (traficAlert *TraficAlert) AverageTrafic() int64 {
	return traficAlert.averageTrafic
}

func (traficAlert *TraficAlert) AboveThreshold() bool {
	return traficAlert.direction
}

func (traficAlert *TraficAlert) UnderThreshold() bool {
	return !traficAlert.direction
}

func newTraficAlert(triggeredAt time.Time, alert string, averageTrafic int64, direction bool) TraficAlert {
	return TraficAlert{
		triggeredAt:   triggeredAt,
		alert:         alert,
		averageTrafic: averageTrafic,
		direction:     direction,
	}
}

func newAlertMonitor(period time.Duration, threshold int64) alertMonitor {
	return alertMonitor{
		period:            period,
		logList:           list.New(),
		threshold:         threshold,
		totalTrafic:       0,
		wasAboveThreshold: false,
	}
}

func (alertMonitor *alertMonitor) addLogEntry(logEntry httplog.LogEntry) httplog.Alert {
	alertMonitor.logList.PushFront(logEntry)
	alertMonitor.totalTrafic++

	alertMonitor.invalidateLogsBefore(logEntry.Time.Add(-alertMonitor.period))

	alert := alertMonitor.checkTrafic(logEntry.Time)
	return alert
}

func (alertMonitor *alertMonitor) invalidateLogsBefore(time time.Time) {
	l := alertMonitor.logList
	for entry := l.Front(); entry != nil && entry.Value.(httplog.LogEntry).Time.Before(time); entry = entry.Next() {
		alertMonitor.totalTrafic--
		l.Remove(entry)
	}
}

func (alertMonitor *alertMonitor) isAboveThreshold() bool {
	return alertMonitor.totalTrafic > alertMonitor.threshold
}

func (alertMonitor *alertMonitor) checkTrafic(at time.Time) httplog.Alert {
	if alertMonitor.isAboveThreshold() {
		if alertMonitor.wasAboveThreshold {
			alertMonitor.wasAboveThreshold = false
			return newTraficAlert(at, "Traffic is back to normal - Hits: %v, Triggered at %v", alertMonitor.totalTrafic, false)
		}
		alertMonitor.wasAboveThreshold = true
		return newTraficAlert(at, "High traffic generated an alert - Hits: %v, Triggered at %v", alertMonitor.totalTrafic, true)
	}
	return nil
}
