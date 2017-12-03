package httplog

import "time"

type Alert interface {
	TriggeredAt() time.Time
	Alert() string
}
