package main

import (
	ui "github.com/gizak/termui"
)

type CPUUsagePercent struct {
	Pct  float64
	Data []float64
}

func NewCpuUsageWidget() *ui.Par {
	cpuUsage := ui.NewPar("...Awaiting CPU Stats")
	cpuUsage.TextFgColor = ui.ColorWhite
	cpuUsage.Height = 5
	cpuUsage.BorderLabel = "CPU Usage"
	cpuUsage.BorderFg = ui.ColorCyan
	return cpuUsage
}

func NewCpuUsageChart() *ui.LineChart {
	cpuGraph := ui.NewLineChart()
	cpuGraph.BorderLabel = "CPU Usage"
	cpuGraph.Height = 10
	cpuGraph.X = 0
	cpuGraph.Y = 0
	cpuGraph.AxesColor = ui.ColorWhite
	cpuGraph.LineColor = ui.ColorRed
	return cpuGraph
}
