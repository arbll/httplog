package main

import (
	"flag"
	"log"
	"time"

	ui "github.com/gizak/termui"
	"github.com/omen-/httplog/pkg/commonformat"
	"github.com/omen-/httplog/pkg/monitor"
)

func main() {
	help := flag.Bool("h", false, "Show usage")
	alertThreshold := flag.Int64("threshold", 150, "Traffic threshold after which an alert will be generated.")
	alertPeriod := flag.Duration("aperiod", 2*time.Minute, "An alert will be generated if the trafic for the past <aperiod> minutes exceed the given threshold.")
	trafficReportPeriod := flag.Duration("rperiod", 10*time.Second, "Frequency at which reports will be generated.")
	logPath := flag.String("logpath", "access.log", "Path to the common log format file.")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	reader, err := commonformat.NewReader(*logPath, &commonformat.LogParser{})
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	mui := buildUI()

	config := monitor.Config{
		TrafficReportPeriod: *trafficReportPeriod,
		AlertPeriod:         *alertPeriod,
		AlertThreshold:      *alertThreshold,
	}

	go monitorFile(config, reader, mui)

	ui.Loop()
}

func monitorFile(config monitor.Config, reader *commonformat.Reader, mui *monitorUI) {

	monitor := monitor.New(config, reader)

	for {
		select {
		case logReport := <-monitor.TrafficReports():
			mui.displayReport(logReport)
		case alert := <-monitor.Alerts():
			mui.displayAlert(alert)
		}
	}
}
