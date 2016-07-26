package main

import (
	"fmt"
	"github.com/docker/engine-api/types"
	ui "github.com/gizak/termui"
)

type MemoryUsageWidget struct {
	Views   []ui.GridBufferer
	Handler func(ui.Event)
}

func NewMemoryUsageWidget() *MemoryUsageWidget {
	memoryUsage := ui.NewGauge()
	memoryUsage.Label = "{{percent}}%"
	memoryUsage.Height = 4
	memoryUsage.BorderLabel = "Memory Usage"
	memoryUsage.BarColor = ui.ColorGreen
	memoryUsage.BorderFg = ui.ColorCyan
	//memoryUsage.BorderLabelFg = ui.ColorMagenta

	maxMemoryUsage := ui.NewPar("...Awaiting Memory Stats")
	maxMemoryUsage.TextFgColor = ui.ColorWhite
	maxMemoryUsage.Height = 5
	maxMemoryUsage.BorderLabel = "Max Memory Usage"
	maxMemoryUsage.BorderFg = ui.ColorCyan

	return &MemoryUsageWidget{Views: []ui.GridBufferer{memoryUsage, maxMemoryUsage}, Handler: func(e ui.Event) {
		stats := e.Data.(types.StatsJSON)
		memoryUsage.BorderLabel = fmt.Sprintf("Memory Usage: %d / %d", stats.MemoryStats.Usage, stats.MemoryStats.Limit)
		memoryUsage.Percent = int((float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit)) * 100)
		maxMemoryUsage.Text = fmt.Sprintf("Max Memory Usage: %d", stats.MemoryStats.MaxUsage)

	}}
}
