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
	memoryUsage := ui.NewPar("...Awaiting Memory Stats")
	memoryUsage.TextFgColor = ui.ColorWhite
	memoryUsage.Height = 5
	memoryUsage.BorderLabel = "Memory Usage"
	memoryUsage.BorderFg = ui.ColorCyan

	maxMemoryUsage := ui.NewPar("...Awaiting Memory Stats")
	maxMemoryUsage.TextFgColor = ui.ColorWhite
	maxMemoryUsage.Height = 5
	maxMemoryUsage.BorderLabel = "Max Memory Usage"
	maxMemoryUsage.BorderFg = ui.ColorCyan

	return &MemoryUsageWidget{Views: []ui.GridBufferer{memoryUsage, maxMemoryUsage}, Handler: func(e ui.Event) {
		stats := e.Data.(types.StatsJSON)
		memoryUsage.Text = fmt.Sprintf("Memory Usage: %d / %d", stats.MemoryStats.Usage, stats.MemoryStats.Limit)
		maxMemoryUsage.Text = fmt.Sprintf("Max Memory Usage: %d", stats.MemoryStats.MaxUsage)

	}}
}
