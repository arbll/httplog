package httplog

type LogStream interface {
	Logs() chan LogEntry
}
