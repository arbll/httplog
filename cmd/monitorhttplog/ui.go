package main

import (
	"fmt"
	"sort"

	ui "github.com/gizak/termui"
	"github.com/omen-/httplog/pkg/monitor"
)

const sparklinesBufferSize = 160

type monitorUI struct {
	requestSparklines *ui.Sparklines
	dataSparklines    *ui.Sparklines
	sectionList       *ui.List
	alertList         *ui.List
	alerts            []string
	alertOffset       int
}

func buildUI() *monitorUI {

	splReq := ui.NewSparkline()
	splReq.Height = 8
	splReq.LineColor = ui.ColorCyan
	splReq.Data = []int{}

	splsReq := ui.NewSparklines(splReq)
	splsReq.Height = 11
	splsReq.BorderFg = ui.ColorWhite
	splsReq.BorderLabel = "Requests / 10s"

	splData := ui.NewSparkline()
	splData.Height = 8
	splData.LineColor = ui.ColorRed
	splData.Data = []int{}

	splsData := ui.NewSparklines(splData)
	splsData.Height = 11
	splsData.BorderFg = ui.ColorWhite
	splsData.BorderLabel = "Bytes sent / 10s"

	lsSections := ui.NewList()
	lsSections.Items = []string{}
	lsSections.ItemFgColor = ui.ColorYellow
	lsSections.BorderLabel = "Req / Section (10s)"
	lsSections.Height = 10

	lsAlerts := ui.NewList()
	lsAlerts.Items = []string{}
	lsAlerts.ItemFgColor = ui.ColorRed
	lsAlerts.BorderLabel = "Alerts"
	lsAlerts.Height = 10

	parTips := ui.NewPar("<up>/<down>: Alert history | c: Clear alerts | q: Quit")
	parTips.Height = 3
	parTips.BorderLabel = "Usage"

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, splsReq), ui.NewCol(6, 0, splsData)),
		ui.NewRow(
			ui.NewCol(4, 0, lsSections), ui.NewCol(8, 0, lsAlerts)),
		ui.NewRow(
			ui.NewCol(12, 0, parTips)))

	mui := &monitorUI{
		requestSparklines: splsReq,
		dataSparklines:    splsData,
		sectionList:       lsSections,
		alertList:         lsAlerts,
		alerts:            []string{},
	}

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		if mui.alertOffset > 0 {
			mui.alertOffset--
			mui.refreshAlerts()
		}
	})
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		if mui.alertOffset < len(mui.alerts)-(mui.alertList.Height-2) {
			mui.alertOffset++
			mui.refreshAlerts()
		}
	})
	ui.Handle("/sys/kbd/c", func(ui.Event) {
		mui.alertOffset = 0
		mui.alerts = []string{}
		mui.refreshAlerts()
	})

	ui.Handle("/sys/wnd/resize", func(ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Body.Align()
	ui.Render(ui.Body)

	return mui
}

func (mui *monitorUI) displayReport(report monitor.TraficReport) {
	mui.requestSparklines.Lines[0].Data = append(mui.requestSparklines.Lines[0].Data, int(report.RequestCount))
	mui.dataSparklines.Lines[0].Data = append(mui.dataSparklines.Lines[0].Data, int(report.BytesSentCount))

	sections := make([]string, 0, len(report.RequestsBySection))

	type secReq struct {
		sec string
		req int64
	}

	var sortedSecs []secReq
	for k, v := range report.RequestsBySection {
		sortedSecs = append(sortedSecs, secReq{k, v})
	}

	sort.Slice(sortedSecs, func(i, j int) bool {
		return sortedSecs[i].req > sortedSecs[j].req
	})

	for _, secReq := range sortedSecs {
		sections = append(sections, fmt.Sprintf("%v %v", secReq.req, secReq.sec))
	}
	mui.sectionList.Items = sections

	ui.Render(ui.Body)
}

func (mui *monitorUI) displayAlert(alert monitor.Alert) {
	mui.alerts = append(mui.alerts, alert.Alert())
	mui.refreshAlerts()
}

func (mui *monitorUI) refreshAlerts() {
	space := mui.alertList.GetHeight() - 2

	if space >= len(mui.alerts) {
		mui.alertList.Items = mui.alerts
	} else {
		mui.alertList.Items = mui.alerts[mui.alertOffset:]
	}

	ui.Render(ui.Body)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
