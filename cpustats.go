package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

type CpuUsageWidget struct {
	*ui.Par
	Handler func(ui.Event)
}

type CpuChart struct {
	*ui.LineChart
	Handler func(e ui.Event)
}

type CPUUsagePercent struct {
	Pct  float64
	Data []float64
}

func NewCpuUsageWidget() *CpuUsageWidget {
	cpuUsage := ui.NewPar("...Awaiting CPU Stats")
	cpuUsage.TextFgColor = ui.ColorWhite
	cpuUsage.Height = 5
	cpuUsage.BorderLabel = "CPU Usage"
	cpuUsage.BorderFg = ui.ColorCyan
	return &CpuUsageWidget{Par: cpuUsage, Handler: func(e ui.Event) {
		stats := e.Data.(CPUUsagePercent)
		cpuUsage.Text = fmt.Sprintf("CPU Usage: %5.2f%%", stats.Pct*100)

	}}
}

func NewCpuUsageChart() *CpuChart {
	cpuGraph := ui.NewLineChart()
	cpuGraph.BorderLabel = "CPU Usage"
	cpuGraph.Height = 10
	cpuGraph.X = 0
	cpuGraph.Y = 0
	cpuGraph.AxesColor = ui.ColorWhite
	cpuGraph.LineColor = ui.ColorRed
	return &CpuChart{LineChart: cpuGraph, Handler: func(e ui.Event) {
		stats := e.Data.(CPUUsagePercent)
		cpuGraph.Data = stats.Data
	}}
}
