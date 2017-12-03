package monitor

import (
	"container/list"
	"fmt"
	"time"

	"github.com/omen-/httplog/pkg/commonformat"
)

type alertMonitor struct {
	period            time.Duration
	logList           *list.List
	totalTrafic       int64
	threshold         int64
	wasAboveThreshold bool
}

// Alert is an interface used to represent an alert raised by the monitor.
type Alert interface {
	// TriggeredAt returns the instant at which the alert was triggered.
	TriggeredAt() time.Time
	// Alert returns the formated alert text.
	Alert() string
}

// TraficAlert represents a trafic alert.
type TraficAlert struct {
	Triggered     time.Time
	AlertText     string
	AverageTrafic int64
	direction     bool
}

func newTraficAlert(triggeredAt time.Time, alert string, averageTrafic int64, direction bool) *TraficAlert {
	return &TraficAlert{
		Triggered:     triggeredAt,
		AlertText:     alert,
		AverageTrafic: averageTrafic,
		direction:     direction,
	}
}

// Alert returns the formated alert text.
func (ta *TraficAlert) Alert() string {
	return ta.AlertText
}

// TriggeredAt returns the instant at which the alert was triggered.
func (ta *TraficAlert) TriggeredAt() time.Time {
	return ta.Triggered
}

// AboveThreshold returns true if the alert was triggered because the trafic
// went past the threshold.
func (ta *TraficAlert) AboveThreshold() bool {
	return ta.direction
}

// UnderThreshold returns true if the alert was triggered because the trafic
// went back under the threshold.
func (ta *TraficAlert) UnderThreshold() bool {
	return !ta.direction
}

// newAlertMonitor returns a new alert monitor that will generate alerts when
// the trafic during the given period exceed the given threshold.
func newAlertMonitor(period time.Duration, threshold int64) *alertMonitor {
	return &alertMonitor{
		period:            period,
		logList:           list.New(),
		threshold:         threshold,
		totalTrafic:       0,
		wasAboveThreshold: false,
	}
}

// addLogEntry adds a new log entry to the log buffer.
// Will raise an Alert if the new entry makes the traffic go above or under the
// monitor threshold.
func (alertMonitor *alertMonitor) addLogEntry(logEntry commonformat.LogEntry) Alert {
	alertMonitor.logList.PushBack(logEntry)
	alertMonitor.totalTrafic++

	alertMonitor.invalidateLogsBefore(logEntry.Time.Add(-alertMonitor.period))

	alert := alertMonitor.checkTrafic(logEntry.Time)
	return alert
}

// invalidateLogsBefore removes all the log entries before the given time.
func (alertMonitor *alertMonitor) invalidateLogsBefore(time time.Time) {
	l := alertMonitor.logList
	var next *list.Element
	for entry := l.Front(); entry != nil && entry.Value.(commonformat.LogEntry).Time.Before(time); entry = next {
		next = entry.Next()
		alertMonitor.totalTrafic--
		l.Remove(entry)
	}
}

func (alertMonitor *alertMonitor) isAboveThreshold() bool {
	return alertMonitor.totalTrafic > alertMonitor.threshold
}

// checkTrafic checks if the current state of the trafic should raise an Alert.
// The given instant will be used to set the Alert trigger time. time.Now()
// should not be used since the logs might be delayed.
func (alertMonitor *alertMonitor) checkTrafic(at time.Time) Alert {
	if alertMonitor.isAboveThreshold() && !alertMonitor.wasAboveThreshold {
		alertMonitor.wasAboveThreshold = true
		alert := fmt.Sprintf("[%v] High traffic generated an alert - Hits: %v", at.Format("15:04"), alertMonitor.totalTrafic)
		return newTraficAlert(at, alert, alertMonitor.totalTrafic, true)
	} else if !alertMonitor.isAboveThreshold() && alertMonitor.wasAboveThreshold {
		alertMonitor.wasAboveThreshold = false
		alert := fmt.Sprintf("[%v] Traffic is back to normal - Hits: %v", at.Format("15:04"), alertMonitor.totalTrafic)
		return newTraficAlert(at, alert, alertMonitor.totalTrafic, false)
	}
	return nil
}
