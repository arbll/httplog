package main

import (
	"log"

	ui "github.com/gizak/termui"
	"github.com/omen-/httplog/pkg/commonformat"
	"github.com/omen-/httplog/pkg/monitor"
)

func main() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	reader, err := commonformat.NewReader("access.log", commonformat.LogParser{})
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	mui := buildUI()

	go monitorFile(reader, mui)

	ui.Loop()
}

func monitorFile(reader *commonformat.Reader, mui *monitorUI) {

	monitor := monitor.New(10, reader)

	for {
		select {
		case logReport := <-monitor.TrafficReports():
			mui.displayReport(logReport)
		case alert := <-monitor.Alerts():
			mui.displayAlert(alert)
		}
	}
}
