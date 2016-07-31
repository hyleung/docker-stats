package main

import (
	"fmt"
	//log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	humanize "github.com/dustin/go-humanize"
	ui "github.com/gizak/termui"
)

type BlkIOWidget struct {
	Views   []ui.GridBufferer
	Handler func(ui.Event)
}

func NewBlkIOWidget() BlkIOWidget {
	servicedOps := ui.NewList()
	servicedOps.BorderLabel = "I/O Operations"
	servicedOps.BorderFg = ui.ColorCyan
	servicedOps.PaddingLeft = 1
	servicedBytes := ui.NewList()
	servicedBytes.BorderLabel = "I/O Bytes"
	servicedBytes.BorderFg = ui.ColorCyan
	servicedBytes.PaddingLeft = 1

	servicedOps.Height = 5
	servicedBytes.Height = 5
	return BlkIOWidget{Views: []ui.GridBufferer{servicedOps, servicedBytes}, Handler: func(e ui.Event) {
		stats := e.Data.(types.StatsJSON)
		opsData := make([]string, len(stats.BlkioStats.IoServicedRecursive))
		for idx, element := range stats.BlkioStats.IoServicedRecursive {
			opsData[idx] = fmt.Sprintf("%s: %d", element.Op, element.Value)
		}
		servicedOps.Items = opsData
		bytesData := make([]string, len(stats.BlkioStats.IoServiceBytesRecursive))
		for idx, element := range stats.BlkioStats.IoServiceBytesRecursive {
			bytesData[idx] = fmt.Sprintf("%s: %s", element.Op, humanize.Bytes(element.Value))
		}
		servicedBytes.Items = bytesData
	}}
}
