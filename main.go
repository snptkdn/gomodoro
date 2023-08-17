package main

import (
	"fmt"
	"time"

	"github.com/rivo/tview"
)

const refreshInterval = 1000 * time.Millisecond

var (
	view  *tview.Modal
	app   *tview.Application
	timer time.Time
)

func timeString() string {
	return fmt.Sprintf(timer.Format("15:04:05"))
}

func tick() {
	timer = timer.Add(time.Second)
}

func updateTime() {
	for {
		time.Sleep(refreshInterval)
		tick()
		app.QueueUpdateDraw(func() {
      if 
			view.SetText(timeString())
		})
	}
}

func main() {
	timer = time.Time{}
	app = tview.NewApplication()
  view = tview.NewModal().SetText(timeString()).AddButtons([]string{"▶︎","■"})

	go updateTime()
	if err := app.SetRoot(view, false).Run(); err != nil {
		panic(err)
	}
}
