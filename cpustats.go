package main

import (
	"fmt"
	"github.com/docker/engine-api/types"
	ui "github.com/gizak/termui"
)

type CpuUsageWidget struct {
	*ui.Par
	*ui.LineChart
	Handler func(ui.Event)
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

	cpuGraph := ui.NewLineChart()
	cpuGraph.BorderLabel = "CPU Usage"
	cpuGraph.Height = 10
	cpuGraph.PaddingTop = 1
	cpuGraph.Mode = "dot"
	cpuGraph.X = 0
	cpuGraph.Y = 0
	cpuGraph.AxesColor = ui.ColorWhite
	cpuGraph.LineColor = ui.ColorRed
	var i = 0
	var currentCPUUsage = uint64(0)
	var currentSystemUsage = uint64(0)
	var cpuHistory = make([]float64, 600)
	var cpuHead = 0
	return &CpuUsageWidget{Par: cpuUsage, LineChart: cpuGraph, Handler: func(e ui.Event) {
		stats := e.Data.(types.StatsJSON)
		var cpuPct = 0.0
		cpuPct, currentCPUUsage, currentSystemUsage = computeCpu(stats, currentCPUUsage, currentSystemUsage)
		cpuHistory[cpuHead] = cpuPct * 100.0
		if cpuHead < 599 {
			cpuHead = cpuHead + 1
		} else {
			cpuHead = 0
		}
		cpuUsage.Text = fmt.Sprintf("CPU Usage: %5.2f%%", cpuPct*100)
		cpuGraph.Data = cpuHistory[:i]
		i = i + 1
	}}
}

func computeCpu(stats types.StatsJSON, currentUsage uint64, currentSystemUsage uint64) (cpuPct float64, cpuUsage uint64, systemUsage uint64) {
	//compute the cpu usage percentage
	//via https://github.com/docker/docker/blob/e884a515e96201d4027a6c9c1b4fa884fc2d21a3/api/client/container/stats_helpers.go#L199-L212
	newCpuUsage := stats.CPUStats.CPUUsage.TotalUsage
	newSystemUsage := stats.CPUStats.SystemUsage
	cpuDiff := float64(newCpuUsage) - float64(currentUsage)
	systemDiff := float64(newSystemUsage) - float64(currentSystemUsage)
	return cpuDiff / systemDiff * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)), newCpuUsage, newSystemUsage
}
