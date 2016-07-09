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
		p.Height = 3
		p.Width = 50
		p.TextFgColor = ui.ColorWhite
		p.BorderLabel = "Text Box"
		p.BorderFg = ui.ColorCyan

		cpu := ui.NewPar("...Awaiting CPU Stats")
		cpu.Height = 3
		cpu.Width = 50
		cpu.Y = 5
		cpu.TextFgColor = ui.ColorWhite
		cpu.BorderLabel = "Container CPU"
		cpu.BorderFg = ui.ColorCyan

		g := ui.NewGauge()
		g.Percent = 0
		g.Width = 50
		g.Height = 3
		g.Y = 11
		g.BorderLabel = "Guage"
		g.BarColor = ui.ColorRed
		g.BorderFg = ui.ColorWhite
		g.BorderLabelFg = ui.ColorCyan

		ui.Render(p, g, cpu)

		ui.Handle("/sys/kbd/q", func(ui.Event) {
			ui.StopLoop()
		})
		ui.Handle("/docker/stats", func(e ui.Event) {
			stats := e.Data.(types.Stats)
			cpu := ui.NewPar(fmt.Sprintf("CPU Usage: %d", stats.CPUStats.SystemUsage))
			cpu.Height = 3
			cpu.Width = 50
			cpu.Y = 5
			cpu.TextFgColor = ui.ColorWhite
			cpu.BorderLabel = "Container CPU"
			cpu.BorderFg = ui.ColorCyan
			ui.Render(cpu)
		})
		ui.Loop()
		return nil
	}
	app.Run(os.Args)
}
