package main

import (
	"fmt"
	"github.com/docker/engine-api/types"
	ui "github.com/gizak/termui"
)

type CpuUsageWidget struct {
	Views   []ui.GridBufferer
	Handler func(ui.Event)
}

type CPUUsagePercent struct {
	Pct  float64
	Data []float64
}

func NewCpuUsageWidget() *CpuUsageWidget {
	cpuGraph := ui.NewLineChart()
	cpuGraph.BorderLabel = "CPU Usage"
	cpuGraph.Height = 10
	cpuGraph.PaddingTop = 1
	cpuGraph.PaddingRight = 10
	cpuGraph.Mode = "dot"
	cpuGraph.X = 0
	cpuGraph.Y = 0
	cpuGraph.AxesColor = ui.ColorWhite
	cpuGraph.LineColor = ui.ColorRed
	cpuGraph.BorderFg = ui.ColorCyan
	var i = 0
	var currentCPUUsage = uint64(0)
	var currentSystemUsage = uint64(0)
	var cpuHistory = make([]float64, 600)
	var cpuHead = 0
	return &CpuUsageWidget{Views: []ui.GridBufferer{cpuGraph}, Handler: func(e ui.Event) {
		stats := e.Data.(types.StatsJSON)
		var cpuPct = 0.0
		cpuPct, currentCPUUsage, currentSystemUsage = computeCpu(stats, currentCPUUsage, currentSystemUsage)
		cpuHistory[cpuHead] = cpuPct * 100.0
		if cpuHead < len(cpuHistory) {
			cpuHead = cpuHead + 1
		} else {
			cpuHead = 0
			//reset the data range
			cpuHistory = make([]float64, 600)
		}
		cpuGraph.BorderLabel = fmt.Sprintf("CPU Usage: %5.2f%%", cpuPct*100)
		cpuGraph.Data = getDataRange(cpuGraph, cpuHistory, cpuHead)
		i = i + 1
	}}
}

func computeNumPoints(lc *ui.LineChart) int {
	padding := 9 * 2
	return (lc.Width - padding)
}
func getDataRange(lc *ui.LineChart, data []float64, head int) []float64 {
	points := computeNumPoints(lc)
	if head < points {
		return data
	}
	return data[(head - points):]
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
