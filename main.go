package main

import (
	"bufio"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	ui "github.com/gizak/termui"
	"golang.org/x/net/context"
	"os"
	"strings"
)

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		containerName := c.Args().Get(0)
		log.Info("Starting monitoring on ", containerName)
		go func() {
			//start watching the stats feed
			defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
			cli, err := client.NewClient("unix:///var/run/docker.sock", "1.24", nil, defaultHeaders)
			if err != nil {
				panic(err)
			}
			readCloser, err := cli.ContainerStats(context.Background(), containerName, true)
			defer readCloser.Close()
			if err != nil {
				panic(err)
			}
			scanner := bufio.NewScanner(readCloser)
			for scanner.Scan() {
				var stats types.StatsJSON
				err = json.NewDecoder(strings.NewReader(scanner.Text())).Decode(&stats)
				if err != nil {
					panic(err)
				}
				ui.SendCustomEvt("/docker/stats", stats)
			}
		}()
		err := ui.Init()
		if err != nil {
			panic(err)
		}
		defer ui.Close()
		p := ui.NewPar(":PRESS q to QUIT")
		p.TextFgColor = ui.ColorWhite
		p.BorderLabel = "Menu"
		p.Height = 3
		p.BorderFg = ui.ColorCyan

		cpuUsage := NewCpuUsageWidget()
		memoryUsage := NewMemoryUsageWidget()
		networkStats := NewNetworkStats()
		blkIOStats := NewBlkIOWidget()
		//Grid layout
		ui.Body.AddRows(
			ui.NewRow(
				ui.NewCol(12, 0, cpuUsage.Views...),
			),
			ui.NewRow(
				ui.NewCol(3, 0, networkStats.RxViews...),
				ui.NewCol(3, 0, networkStats.TxViews...),
				ui.NewCol(6, 0, memoryUsage.Views...),
			),
			ui.NewRow(
				ui.NewCol(3, 0, blkIOStats.Views[0]),
				ui.NewCol(3, 0, blkIOStats.Views[1]),
			),
			ui.NewRow(
				ui.NewCol(12, 0, p),
			),
		)
		ui.Body.Align()
		ui.Render(ui.Body)

		ui.Handle("/sys/kbd/q", func(ui.Event) {
			ui.StopLoop()
		})
		ui.Handle("/docker/stats", func(e ui.Event) {
			networkStats.Handler(e)
			memoryUsage.Handler(e)
			cpuUsage.Handler(e)
			blkIOStats.Handler(e)
			ui.Render(ui.Body)
		})
		ui.Loop()
		return nil
	}
	app.Run(os.Args)
}
