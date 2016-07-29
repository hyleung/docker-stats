package main

import (
	"fmt"
	//log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	ui "github.com/gizak/termui"
)

type BlkIOWidget struct {
	Views   []ui.GridBufferer
	Handler func(ui.Event)
}

func NewBlkIOWidget() BlkIOWidget {
	servicedOps := ui.NewList()
	servicedOps.BorderLabel = "Service Operations"
	servicedOps.BorderFg = ui.ColorCyan

	servicedOps.Height = 5
	return BlkIOWidget{Views: []ui.GridBufferer{servicedOps}, Handler: func(e ui.Event) {
		stats := e.Data.(types.StatsJSON)
		data := make([]string, len(stats.BlkioStats.IoServicedRecursive))
		for idx, element := range stats.BlkioStats.IoServicedRecursive {
			data[idx] = fmt.Sprintf("%s: %d", element.Op, element.Value)
		}
		servicedOps.Items = data
	}}
}
