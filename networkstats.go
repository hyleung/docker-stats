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

	rxData := []string{
		"RxBytes:   0",
		"RxPackets: 0",
		"RxErrors:  0",
		"RxDropped: 0",
	}

	rxList.Items = rxData
	rxList.Height = 10
	rxList.PaddingLeft = 1
	rxViews := make([]ui.GridBufferer, 4)
	rxViews[0] = rxList

	txList := ui.NewList()
	txList.BorderLabel = "Network Transmit"

	txData := []string{
		"TxBytes:   0",
		"TxPackets: 0",
		"TxErrors:  0",
		"TxDropped: 0",
	}

	txList.Items = txData
	txList.Height = 10
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
		fmt.Sprintf("RxBytes:   %d", stats.RxBytes),
		fmt.Sprintf("RxPackets: %d", stats.RxPackets),
		fmt.Sprintf("RxErrors:  %d", stats.RxErrors),
		fmt.Sprintf("RxDropped: %d", stats.RxDropped),
	}
	return data
}
func formatTxData(stats types.NetworkStats) []string {
	data := []string{
		fmt.Sprintf("TxBytes:   %d", stats.TxBytes),
		fmt.Sprintf("TxPackets: %d", stats.TxPackets),
		fmt.Sprintf("TxErrors:  %d", stats.TxErrors),
		fmt.Sprintf("TxDropped: %d", stats.TxDropped),
	}
	return data
}
