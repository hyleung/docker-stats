package main

import (
	//	"fmt"
	//	"github.com/docker/engine-api/types"
	//	humanize "github.com/dustin/go-humanize"
	ui "github.com/gizak/termui"
)

type MenuWidget struct {
	ui.GridBufferer
}

func NewMenu() *MenuWidget {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	p := ui.NewList()
	p.Items = []string{
		"Press 'q' to quit",
	}
	p.BorderLabel = "Menu"
	p.Height = 3
	return &MenuWidget{p}
}
