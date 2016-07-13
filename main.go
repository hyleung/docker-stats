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

type CPUUsagePercent struct {
	Pct  float64
	Data []float64
}

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
			var currentCPUUsage = uint64(0)
			var currentSystemUsage = uint64(0)
			var cpuHistory = make([]float64, 60)
			var cpuHead = 0
			for scanner.Scan() {
				var stats types.Stats
				err = json.NewDecoder(strings.NewReader(scanner.Text())).Decode(&stats)
				if err != nil {
					panic(err)
				}
				ui.SendCustomEvt("/docker/stats", stats)
				//compute the cpu usage percentage
				//via https://github.com/docker/docker/blob/e884a515e96201d4027a6c9c1b4fa884fc2d21a3/api/client/container/stats_helpers.go#L199-L212

				cpuDiff := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(currentCPUUsage)
				systemDiff := float64(stats.CPUStats.SystemUsage) - float64(currentSystemUsage)
				cpuPct := cpuDiff / systemDiff * float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
				cpuHistory[cpuHead] = cpuPct * 100.0
				if cpuHead < 59 {
					cpuHead = cpuHead + 1
				} else {
					cpuHead = 0
				}
				ui.SendCustomEvt("/docker/cpuPct", CPUUsagePercent{Pct: cpuPct, Data: cpuHistory})
				currentCPUUsage = stats.CPUStats.CPUUsage.TotalUsage
				currentSystemUsage = stats.CPUStats.SystemUsage
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

		cpuUsage := ui.NewPar("...Awaiting CPU Stats")
		cpuUsage.TextFgColor = ui.ColorWhite
		cpuUsage.Height = 5
		cpuUsage.BorderLabel = "CPU Usage"
		cpuUsage.BorderFg = ui.ColorCyan

		cpuGraph := ui.NewLineChart()
		cpuGraph.BorderLabel = "CPU Usage"
		cpuGraph.Height = 10
		cpuGraph.X = 0
		cpuGraph.Y = 0
		cpuGraph.AxesColor = ui.ColorWhite
		cpuGraph.LineColor = ui.ColorRed

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
				ui.NewCol(3, 0, cpuUsage),
				ui.NewCol(3, 0, cpuGraph),
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
			memoryUsage.Text = fmt.Sprintf("Memory Usage: %d / %d", stats.MemoryStats.Usage, stats.MemoryStats.Limit)
			maxMemoryUsage.Text = fmt.Sprintf("Max Memory Usage: %d", stats.MemoryStats.MaxUsage)
			ui.Render(ui.Body)
		})
		ui.Handle("/docker/cpuPct", func(e ui.Event) {
			stats := e.Data.(CPUUsagePercent)
			cpuUsage.Text = fmt.Sprintf("CPU Usage: %5.2f%%", stats.Pct*100)
			cpuGraph.Data = stats.Data
		})
		ui.Loop()
		return nil
	}
	app.Run(os.Args)
}
