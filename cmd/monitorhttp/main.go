package main

import (
	"log"

	"github.com/omen-/httplog"

	ui "github.com/gizak/termui"
	"github.com/omen-/httplog/commonformat"
	"github.com/omen-/httplog/logfile"
	"github.com/omen-/httplog/logmonitor"
)

func main() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	reader, err := logfile.NewReader("access.log", commonformat.LogParser{})
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	mui := buildUI()

	go monitor(reader, mui)

	ui.Loop()
}

func monitor(reader httplog.LogStream, mui *monitorUI) {

	monitor := logmonitor.New(10, reader)

	for {
		select {
		case logReport := <-monitor.TraficReports():
			mui.displayReport(logReport)
		case alert := <-monitor.Alerts():
			mui.displayAlert(alert)
		}
	}
}
