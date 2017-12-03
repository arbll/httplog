package main

import (
	"fmt"
	"log"

	"github.com/omen-/httplog/commonformat"
	"github.com/omen-/httplog/logfile"
	"github.com/omen-/httplog/logmonitor"
)

func main() {
	reader, err := logfile.NewReader("access.log", commonformat.LogParser{})
	if err != nil {
		log.Fatalln(err)
	}
	monitor := logmonitor.New(10, reader)

	for {
		select {
		case logReport := <-monitor.TraficReports():
			fmt.Println(logReport)
		case alert := <-monitor.Alerts():
			fmt.Println(alert)
		}
	}
}
