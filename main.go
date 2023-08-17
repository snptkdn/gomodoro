package main

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/rivo/tview"
)

const refreshInterval = 1000 * time.Millisecond

var (
	view            *tview.Modal
	app             *tview.Application
	timer           time.Time
	isPause         = false
	pomodoroEndTime = 5 * time.Second
	breakEndTime    = 10 * time.Second
	lunchEndTime    = 20 * time.Second
	pomodoroCount   = 0
	isBreak         = false
)

func timeString() string {
	return fmt.Sprintf(timer.Format("15:04:05"))
}

func statusString() string {
	if isPause {
    return "PAUSE"
	} else if isBreak && pomodoroCount == 3 {
		return "LUNCH"
	} else if isBreak && pomodoroCount != 4 {
    return fmt.Sprintf("BREAK:%d", pomodoroCount+1)
	} else {
    return fmt.Sprintf("FOCUS:%d", pomodoroCount+1)
	}
}

func tick() {
	timer = timer.Add(time.Second)
}

func beep(status string) {
	err := beeep.Notify("Gomodoro", fmt.Sprintf("%s time is ended!", status), "nothing")
	if err != nil {
		panic(err)
	}
}

func updateTime(ticker time.Ticker, stopTimer chan int) {
	for {
		select {
		case <-ticker.C:
			time.Sleep(refreshInterval)
			tick()
			if !isBreak && timer.Sub(time.Time{}) == pomodoroEndTime {
				beep("Focus")
			}
      if isBreak && pomodoroCount != 4 && timer.Sub(time.Time{}) == breakEndTime {
        beep("Break")
      }
      if isBreak && pomodoroCount == 4 && timer.Sub(time.Time{}) == lunchEndTime {
        beep("Lunch")
      }
			app.QueueUpdateDraw(func() {
				view.SetText(fmt.Sprintf("%s: %s", statusString(), timeString()))
			})
		case <-stopTimer:
			fmt.Println("stopped!")
			return
		}
	}
}

func main() {
	timer = time.Time{}
	ticker := time.NewTicker(1000 * time.Millisecond)
	chStop := make(chan int, 1)
	go updateTime(*ticker, chStop)
	app = tview.NewApplication()
	view = tview.
		NewModal().
		SetText(fmt.Sprintf("%s: %s", statusString(), timeString())).
		AddButtons([]string{"Next", "Pause", "Reset"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Pause" && !isPause {
				ticker.Stop()
				chStop <- 0
				close(chStop)
				isPause = true
			} else if buttonLabel == "Pause" && isPause {
				ticker = time.NewTicker(1000 * time.Millisecond)
				chStop = make(chan int, 1)
				go updateTime(*ticker, chStop)
				isPause = false
			} else if buttonLabel == "Reset" {
				timer = time.Time{}
			} else if buttonLabel == "Next" {
        timer = time.Time{}
        if !isBreak {
          isBreak = true
        } else {
          isBreak = false
          pomodoroCount++
          if pomodoroCount == 4 {
            pomodoroCount = 0
          }
        }
      }
		})

	if err := app.SetRoot(view, false).Run(); err != nil {
		panic(err)
	}
}
