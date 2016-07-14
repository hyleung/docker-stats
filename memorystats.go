package main

import (
	ui "github.com/gizak/termui"
)

func NewMemoryUsagePar() *ui.Par {
	memoryUsage := ui.NewPar("...Awaiting Memory Stats")
	memoryUsage.TextFgColor = ui.ColorWhite
	memoryUsage.Height = 5
	memoryUsage.BorderLabel = "Memory Usage"
	memoryUsage.BorderFg = ui.ColorCyan
	return memoryUsage

}

func NewMaxMemoryWidget() *ui.Par {
	maxMemoryUsage := ui.NewPar("...Awaiting Memory Stats")
	maxMemoryUsage.TextFgColor = ui.ColorWhite
	maxMemoryUsage.Height = 5
	maxMemoryUsage.BorderLabel = "Max Memory Usage"
	maxMemoryUsage.BorderFg = ui.ColorCyan
	return maxMemoryUsage
}
