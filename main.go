package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
				var stats types.Stats
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
		p := ui.NewPar(":PRESS q to QUIT DEMO")
		p.TextFgColor = ui.ColorWhite
		p.BorderLabel = "Text Box"
		p.Height = 3
		p.BorderFg = ui.ColorCyan

		systemUsage := ui.NewPar("...Awaiting CPU Stats")
		systemUsage.TextFgColor = ui.ColorWhite
		systemUsage.Height = 5
		systemUsage.BorderLabel = "System CPU Usage"
		systemUsage.BorderFg = ui.ColorCyan

		cpuUsage := ui.NewPar("...Awaiting CPU Stats")
		cpuUsage.TextFgColor = ui.ColorWhite
		cpuUsage.Height = 5
		cpuUsage.BorderLabel = "CPU Usage"
		cpuUsage.BorderFg = ui.ColorCyan

		memoryUsage := ui.NewPar("...Awaiting Memory Stats")
		memoryUsage.TextFgColor = ui.ColorWhite
		memoryUsage.Height = 5
		memoryUsage.BorderLabel = "Memory Usage"
		memoryUsage.BorderFg = ui.ColorCyan

		maxMemoryUsage := ui.NewPar("...Awaiting Memory Stats")
		maxMemoryUsage.TextFgColor = ui.ColorWhite
		maxMemoryUsage.Height = 5
		maxMemoryUsage.BorderLabel = "Max Memory Usage"
		maxMemoryUsage.BorderFg = ui.ColorCyan
		//Grid layout
		ui.Body.AddRows(
			ui.NewRow(
				ui.NewCol(3, 0, systemUsage),
				ui.NewCol(3, 0, cpuUsage),
				ui.NewCol(3, 0, memoryUsage),
				ui.NewCol(3, 0, maxMemoryUsage),
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
			stats := e.Data.(types.Stats)
			systemUsage.Text = fmt.Sprintf("CPU System Usage: %d", stats.CPUStats.SystemUsage)
			cpuUsage.Text = fmt.Sprintf("CPU Total Usage: %d", stats.CPUStats.CPUUsage.TotalUsage)
			memoryUsage.Text = fmt.Sprintf("Memory Usage: %d / %d", stats.MemoryStats.Usage, stats.MemoryStats.Limit)
			maxMemoryUsage.Text = fmt.Sprintf("Max Memory Usage: %d", stats.MemoryStats.MaxUsage)
			ui.Render(ui.Body)
		})
		ui.Loop()
		return nil
	}
	app.Run(os.Args)
}
