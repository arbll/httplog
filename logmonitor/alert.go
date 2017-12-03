package logmonitor

import (
	"container/list"
	"fmt"
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
	fmt.Println(alertMonitor.logList)
	alertMonitor.logList.PushBack(logEntry)
	alertMonitor.totalTrafic++

	alertMonitor.invalidateLogsBefore(logEntry.Time.Add(-alertMonitor.period))

	alert := alertMonitor.checkTrafic(logEntry.Time)
	return alert
}

func (alertMonitor *alertMonitor) invalidateLogsBefore(time time.Time) {
	l := alertMonitor.logList
	var next *list.Element
	for entry := l.Front(); entry != nil && entry.Value.(httplog.LogEntry).Time.Before(time); entry = next {
		next = entry.Next()
		alertMonitor.totalTrafic--
		l.Remove(entry)
	}
}

func (alertMonitor *alertMonitor) isAboveThreshold() bool {
	return alertMonitor.totalTrafic > alertMonitor.threshold
}

func (alertMonitor *alertMonitor) checkTrafic(at time.Time) httplog.Alert {
	if alertMonitor.isAboveThreshold() && !alertMonitor.wasAboveThreshold {
		alertMonitor.wasAboveThreshold = true
		alert := fmt.Sprintf("High traffic generated an alert - Hits: %v, Triggered at %v", alertMonitor.totalTrafic, at)
		return newTraficAlert(at, alert, alertMonitor.totalTrafic, true)
	} else if !alertMonitor.isAboveThreshold() && alertMonitor.wasAboveThreshold {
		alertMonitor.wasAboveThreshold = false
		alert := fmt.Sprintf("Traffic is back to normal - Hits: %v, Triggered at %v", alertMonitor.totalTrafic, at)
		return newTraficAlert(at, alert, alertMonitor.totalTrafic, false)
	}
	return nil
}
