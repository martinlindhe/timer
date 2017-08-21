package timer

import (
	"fmt"
	"os"
	"time"

	"github.com/martinlindhe/notify"
	"github.com/scryner/systray"
)

// Launch starts the app
func Launch() {
	app := app{
		name: "gotime",
	}

	app.Run()
}

type app struct {
	stopwatchRunning   bool
	stopwatchStartedAt time.Time
	name               string
	menuStopwatch      *systray.MenuItem
}

func (app *app) Run() {
	systray.Run(func() {
		// XXX requires a .ico on windows
		systray.SetIcon(assetData("assets/mac/icon.png"))

		app.menuStopwatch = systray.AddMenuItem("Start stopwatch", "")
		systray.AddMenuSeparatorItem()
		mQuit := systray.AddMenuItem("Quit "+app.name, "")

		// TODO: rather start a gochan when stopwatch starts, and stop it when it ends
		go func() {
			for {
				time.Sleep(500 * time.Millisecond)
				if app.stopwatchRunning {
					app.updateRunningTimer()
				}
			}
		}()

		for {
			select {
			case <-app.menuStopwatch.ClickedCh:
				if !app.stopwatchRunning {
					app.stopwatchRunning = true
					app.stopwatchStartedAt = time.Now()
					app.updateRunningTimer()
					app.menuStopwatch.SetTitle("Stop stopwatch")
					notify.Notify(app.name, "Stopwatch started", "", "assets/icon128.png")
				} else {
					app.stopwatchRunning = false
					systray.SetTitle("")
					app.menuStopwatch.SetTitle("Start stopwatch")
					notify.Notify(app.name, "Stopwatch stopped after "+app.renderRunningTime(), "", "assets/icon128.png")
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	})
}

func (app *app) updateRunningTimer() {
	systray.SetTitle(app.renderRunningTime())
}

func (app *app) renderRunningTime() string {
	duration := time.Now().Sub(app.stopwatchStartedAt)
	s := int(duration.Seconds()) % 60
	m := int(duration.Minutes()) % 60
	h := int(duration.Hours()) % 24
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
