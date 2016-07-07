package main

import (
	"github.com/fatih/color"
	ui "github.com/gizak/termui"
)

func main() {
	color.Red("Hello, world!")
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

	g := ui.NewGauge()
	g.Percent = 0
	g.Width = 50
	g.Height = 3
	g.Y = 11
	g.BorderLabel = "Guage"
	g.BarColor = ui.ColorRed
	g.BorderFg = ui.ColorWhite
	g.BorderLabelFg = ui.ColorCyan

	ui.Render(p, g)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/timer/1s", func(e ui.Event) {
		next := g.Percent + 10
		if next >= 100 {
			g.Percent = 0
		} else {
			g.Percent = next
		}
		ui.Render(g)
	})
	ui.Loop()
}
