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
	systray.Run(onReady)
}

var (
	timerRunning   = false
	timerStartedAt time.Time
	appName        = "gotime"

	mStartTimer *systray.MenuItem
)

func onReady() {
	// XXX requires a .ico on windows
	systray.SetIcon(assetData("assets/mac/icon.png"))

	mStartTimer = systray.AddMenuItem("Start stopwatch", "")
	systray.AddMenuSeparatorItem()
	mQuit := systray.AddMenuItem("Quit "+appName, "")

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			if timerRunning {
				updateRunningTimer()
			}
		}
	}()

	for {
		select {
		case <-mStartTimer.ClickedCh:
			if !timerRunning {
				timerRunning = true
				timerStartedAt = time.Now()
				updateRunningTimer()
				mStartTimer.SetTitle("Stop stopwatch")
				notify.Notify(appName, "Stopwatch started", "", "assets/icon128.png")
			} else {
				timerRunning = false
				systray.SetTitle("")
				mStartTimer.SetTitle("Start stopwatch")
				notify.Notify(appName, "Stopwatch stopped after "+renderRunningTime(), "", "assets/icon128.png")
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
			os.Exit(0)
		}
	}
}

func updateRunningTimer() {
	systray.SetTitle(renderRunningTime())
}

func renderRunningTime() string {
	duration := time.Now().Sub(timerStartedAt)
	s := int(duration.Seconds()) % 60
	m := int(duration.Minutes()) % 60
	h := int(duration.Hours()) % 24
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
