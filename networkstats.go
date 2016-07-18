package main

import (
	//"fmt"
	ui "github.com/gizak/termui"
)

type NetworkStatsWidget struct {
	Views   []ui.GridBufferer
	Handler func(ui.Event)
}

func NewNetworkStats() NetworkStatsWidget {
	list := ui.NewList()
	list.BorderLabel = "Network Stats"

	strs := []string{
		"RxBytes:   0",
		"RxPackets: 0",
		"RxErrors:  0",
		"RxDropped: 0",
		"TxBytes:   0",
		"TxPackets: 0",
		"TxErrors:  0",
		"TxDropped: 0",
	}

	list.Items = strs
	list.Height = 10
	list.PaddingLeft = 1
	result := make([]ui.GridBufferer, 4)
	result[0] = list
	return NetworkStatsWidget{Views: result, Handler: func(e ui.Event) {
		//update the network stats on each event
	}}
}
