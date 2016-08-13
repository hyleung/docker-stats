package main

import (
	"fmt"
	"github.com/docker/engine-api/types"
	//	humanize "github.com/dustin/go-humanize"
	ui "github.com/gizak/termui"
)

type InfoWidget struct {
	ui.GridBufferer
}

func NewInfoWidget(container types.ContainerJSON) *InfoWidget {
	infoList := ui.NewList()
	infoList.BorderLabel = "Container"
	infoList.Items = []string{
		fmt.Sprintf("ID	:	%s", container.ID[:12]),
		fmt.Sprintf("Image:	%s", container.Config.Image),
		//fmt.Sprintf("Cmd:	%s", container.Config.Cmd),
		//fmt.Sprintf("Env:	%s", container.Config.Env),
	}
	infoList.Height = 5
	return &InfoWidget{infoList}
}
