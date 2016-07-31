package main

import (
	"fmt"
	//log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	ui "github.com/gizak/termui"
)

type NetworkStatsWidget struct {
	RxViews []ui.GridBufferer
	TxViews []ui.GridBufferer
	Handler func(ui.Event)
}

func NewNetworkStats() NetworkStatsWidget {
	rxList := ui.NewList()
	rxList.BorderLabel = "Network Received"
	rxList.BorderFg = ui.ColorCyan

	rxData := []string{
		"Bytes:   0",
		"Packets: 0",
		"Errors:  0",
		"Dropped: 0",
	}

	rxList.Items = rxData
	rxList.Height = 6
	rxList.PaddingLeft = 1
	rxViews := make([]ui.GridBufferer, 4)
	rxViews[0] = rxList

	txList := ui.NewList()
	txList.BorderLabel = "Network Transmit"
	txList.BorderFg = ui.ColorCyan
	txData := []string{
		"Bytes:   0",
		"Packets: 0",
		"Errors:  0",
		"Dropped: 0",
	}

	txList.Items = txData
	txList.Height = 6
	txList.PaddingLeft = 1
	txViews := make([]ui.GridBufferer, 4)
	txViews[0] = txList
	return NetworkStatsWidget{RxViews: rxViews, TxViews: txViews, Handler: func(e ui.Event) {
		//update the network stats on each event
		stats := e.Data.(types.StatsJSON)
		for _, value := range stats.Networks {
			rxList.Items = formatRxData(value)
			txList.Items = formatTxData(value)
		}
	}}
}

func formatRxData(stats types.NetworkStats) []string {
	data := []string{
		fmt.Sprintf("Bytes:   %d", stats.RxBytes),
		fmt.Sprintf("Packets: %d", stats.RxPackets),
		fmt.Sprintf("Errors:  %d", stats.RxErrors),
		fmt.Sprintf("Dropped: %d", stats.RxDropped),
	}
	return data
}
func formatTxData(stats types.NetworkStats) []string {
	data := []string{
		fmt.Sprintf("Bytes:   %d", stats.TxBytes),
		fmt.Sprintf("Packets: %d", stats.TxPackets),
		fmt.Sprintf("Errors:  %d", stats.TxErrors),
		fmt.Sprintf("Dropped: %d", stats.TxDropped),
	}
	return data
}
