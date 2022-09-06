package common

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
	"sync"
	"time"
)

var messageColors = []text.Color{
	text.FgRed,
	text.FgGreen,
	text.FgYellow,
	text.FgBlue,
	text.FgMagenta,
	text.FgCyan,
	text.FgWhite,
}

var pw progress.Writer

var tracker *progress.Tracker

func InitProgress() *progress.Writer {
	pw = progress.NewWriter()
	pw.SetAutoStop(false)
	pw.SetTrackerLength(25)
	pw.SetMessageWidth(25)
	//pw.SetNumTrackersExpected(1)
	pw.SetSortBy(progress.SortByPercentDsc)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Options.PercentFormat = "%4.1f%%"
	//pw.Style().Visibility.ETA = true
	//pw.Style().Visibility.ETAOverall = true
	pw.Style().Visibility.Percentage = true
	pw.Style().Visibility.Time = true
	//pw.Style().Visibility.TrackerOverall = true
	pw.Style().Visibility.Value = true
	go pw.Render()

	return &pw
}

func AddTracker(setting *ProgressSetting, wg *sync.WaitGroup) {
	messageFormat := "Current progress %3d"
	tk := &progress.Tracker{Message: messageFormat, Total: setting.Total, Units: *setting.Units}
	tracker = tk

	pw.AppendTracker(tk)
	tick := time.Tick(time.Duration(200) * time.Millisecond)

	wg.Add(1)
	go (func(t *progress.Tracker) {
		defer wg.Done()
		for !t.IsDone() {
			select {

			case <-tick:
				if t.IsDone() {
					print("All Done")
					//renderTable(t, i, rand.Intn(20))
					//tracker.MarkAsDone()
					break
				}
				msg := messageColors[1].Sprint(fmt.Sprintf(messageFormat, t.Value()+1))
				t.UpdateMessage(msg)
			}
		}
	})(tk)

}

func GetTracker() *progress.Tracker {
	return tracker
}
