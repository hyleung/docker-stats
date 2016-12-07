package main

import (
	"bufio"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	ui "github.com/gizak/termui"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"os"
	"strings"
)

const (
	user_agent         = "engine-api-cli-1.0"
	docker_api_version = "1.24"
)

func CreateClient() (*client.Client, error) {
	if os.Getenv("DOCKER_HOST") != "" {
		return client.NewEnvClient()
	} else {
		defaultHeaders := map[string]string{"User-Agent": user_agent}
		return client.NewClient("unix:///var/run/docker.sock", docker_api_version, nil, defaultHeaders)
	}
}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "theme, t",
			Value: "dark",
			Usage: "Theme for UI",
		},
	}

	app.Action = func(c *cli.Context) error {
		containerName := c.Args().Get(0)

		log.Info("Starting monitoring on ", containerName)
		cli, cli_err := CreateClient()
		if cli_err != nil {
			panic(cli_err)
		}
		go func() {
			//start watching the stats feed
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
		if c.String("theme") == "light" {
			ui.ColorMap["border.fg"] = ui.ColorBlack
			ui.ColorMap["label.fg"] = ui.ColorGreen
			ui.ColorMap["linechart.axes.fg"] = ui.ColorBlack
			ui.ColorMap["linechart.line.fg"] = ui.ColorBlack
			ui.ColorMap["gauge.bar.bg"] = ui.ColorGreen
			ui.ColorMap["list.item.fg"] = ui.ColorBlack
			ui.ColorMap["par.text.fg"] = ui.ColorBlack
			ui.ColorMap["gauge.percent.fg"] = ui.ColorBlack
		} else {
			ui.ColorMap["border.fg"] = ui.ColorCyan
			ui.ColorMap["label.fg"] = ui.ColorGreen
			ui.ColorMap["linechart.axes.fg"] = ui.ColorWhite
			ui.ColorMap["linechart.line.fg"] = ui.ColorWhite
			ui.ColorMap["gauge.bar.bg"] = ui.ColorGreen
			ui.ColorMap["list.item.fg"] = ui.ColorWhite
			ui.ColorMap["par.text.fg"] = ui.ColorWhite
			ui.ColorMap["gauge.percent.fg"] = ui.ColorWhite
		}

		//inspect the container
		container, container_err := cli.ContainerInspect(context.Background(), containerName)
		if container_err != nil {
			panic(container_err)
		}
		cpuUsage := NewCpuUsageWidget()
		memoryUsage := NewMemoryUsageWidget()
		networkStats := NewNetworkStats()
		blkIOStats := NewBlkIOWidget()
		menu := NewMenu()
		infoList := NewInfoWidget(container)
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
				ui.NewCol(6, 0, infoList),
			),
			ui.NewRow(
				ui.NewCol(12, 0, menu),
			),
		)
		ui.Body.Align()
		ui.Render(ui.Body)

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
