package timer

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"github.com/martinlindhe/inputbox"
	"github.com/martinlindhe/notify"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

// Launch starts the app
func Launch() {
	app := app{
		name: "gotime",
		icon: "assets/icon128.png",
	}

	app.Run()
}

type app struct {
	stopwatchRunning   bool
	stopwatchStartedAt time.Time
	timerRunning       bool
	timerDuration      time.Duration
	timerStartedAt     time.Time
	name               string
	icon               string
	menuStopwatch      *systray.MenuItem
	menuTimer          *systray.MenuItem
}

func (app *app) Run() {
	systray.Run(func() {
		setTrayIcon()
		app.menuStopwatch = systray.AddMenuItem("Start stopwatch", "")
		app.menuTimer = systray.AddMenuItem("Start timer", "")
		//systray.AddMenuSeparatorItem()
		mQuit := systray.AddMenuItem("Quit "+app.name, "")

		// TODO: rather start a gochan when stopwatch starts, and stop it when it ends
		go func() {
			for {
				time.Sleep(500 * time.Millisecond)
				if app.stopwatchRunning {
					app.updateStopwatch()
				}
				if app.timerRunning {
					app.updateTimer()
				}
			}
		}()

		for {
			select {
			case <-app.menuStopwatch.ClickedCh:
				if !app.stopwatchRunning {
					app.stopwatchRunning = true
					app.stopwatchStartedAt = time.Now()
					app.menuTimer.Disable()
					app.updateStopwatch()
					app.menuStopwatch.SetTitle("Stop stopwatch")
					notify.Notify(app.name, "Stopwatch started", "", app.icon)
				} else {
					app.stopwatchRunning = false
					app.menuTimer.Enable()
					systray.SetTitle("")
					app.menuStopwatch.SetTitle("Start stopwatch")
					notify.Notify(app.name, "Stopwatch stopped after "+app.renderStopwatch(), "", app.icon)
				}
			case <-app.menuTimer.ClickedCh:
				if !app.timerRunning {
					got, ok := inputbox.InputBox("Enter duration", "Enter duration (format: 3h, 5m30s)", "5m")
					if ok {
						var err error
						app.timerDuration, err = time.ParseDuration(got)
						if err != nil {
							notify.Notify(app.name, "Failed to parse duration from input", "", app.icon)
						}

						timer := time.NewTimer(app.timerDuration)
						go func() {
							<-timer.C
							notify.Alert(app.name, "Timer finished after "+renderDuration(app.timerDuration), "", app.icon)
							app.stopTimer()
						}()

						app.timerRunning = true
						app.timerStartedAt = time.Now()
						app.updateTimer()
						app.menuStopwatch.Disable()
						// XXX sound alert when timer is done
						app.menuTimer.SetTitle("Stop timer")
						notify.Notify(app.name, "Timer started", "", app.icon)
					}
				} else {
					app.stopTimer()
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	},
		func() {
			// onExit
		})
}

func (app *app) stopTimer() {
	app.timerRunning = false
	app.menuStopwatch.Enable()
	systray.SetTitle("")
	app.menuTimer.SetTitle("Start timer")
}

func (app *app) updateStopwatch() {
	if runtime.GOOS == "windows" {
		app.menuStopwatch.SetTitle("Stop stopwatch (" + app.renderStopwatch() + ")")
	} else {
		systray.SetTitle(app.renderStopwatch())
	}
}

func (app *app) renderStopwatch() string {
	duration := time.Now().Sub(app.stopwatchStartedAt)
	s := int(duration.Seconds()) % 60
	m := int(duration.Minutes()) % 60
	h := int(duration.Hours()) % 24
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func (app *app) updateTimer() {
	if runtime.GOOS == "windows" {
		app.menuTimer.SetTitle("Stop timer (" + app.renderTimer() + ")")
	} else {
		systray.SetTitle(app.renderTimer())
	}
}

func (app *app) renderTimer() string {
	duration := app.timerStartedAt.Add(app.timerDuration).Sub(time.Now())
	return renderDuration(duration)
}

func renderDuration(dur time.Duration) string {
	s := int(dur.Seconds()) % 60
	m := int(dur.Minutes()) % 60
	h := int(dur.Hours()) % 24
	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm%ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

func setTrayIcon() {
	if runtime.GOOS == "windows" {
		systray.SetIcon(assetData("assets/win/icon.ico"))
	} else {
		systray.SetIcon(assetData("assets/mac/icon.png"))
	}
}
