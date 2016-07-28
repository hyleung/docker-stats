package main

import (
	"fmt"
	"github.com/docker/engine-api/types"
	humanize "github.com/dustin/go-humanize"
	ui "github.com/gizak/termui"
)

type MemoryUsageWidget struct {
	Views   []ui.GridBufferer
	Handler func(ui.Event)
}

func NewMemoryUsageWidget() *MemoryUsageWidget {
	memoryUsage := ui.NewGauge()
	memoryUsage.Label = "{{percent}}%"
	memoryUsage.Height = 3
	memoryUsage.BorderLabel = "Memory Usage"

	memoryUsage.BarColor = ui.ColorGreen
	memoryUsage.BorderFg = ui.ColorCyan

	return &MemoryUsageWidget{Views: []ui.GridBufferer{memoryUsage}, Handler: func(e ui.Event) {
		stats := e.Data.(types.StatsJSON)

		usage := stats.MemoryStats.Usage
		limit := stats.MemoryStats.Limit
		max := stats.MemoryStats.MaxUsage
		memoryUsage.BorderLabel = fmt.Sprintf("Memory Usage: %s / %s (max: %s)", humanize.Bytes(usage), humanize.Bytes(limit), humanize.Bytes(max))
		memoryUsage.Percent = int((float64(stats.MemoryStats.Usage) / float64(stats.MemoryStats.Limit)) * 100)

	}}
}
