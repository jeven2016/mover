package common

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"
	"os"
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

func InitProgress(setting *ProgressSetting, interval time.Duration) {

	pw := progress.NewWriter()
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

	messageFormat := "Downloading Files: %3d"
	tracker := &progress.Tracker{Message: "Progress text", Total: setting.Total, Units: *setting.Units}

	pw.AppendTracker(tracker)

	tick := time.Tick(interval)
	i := 1
	for !tracker.IsDone() {
		select {

		case <-tick:
			tracker.Increment(1)
			if i >= 50 {
				print("All Done")
				//renderTable(t, i, rand.Intn(20))
				tracker.MarkAsDone()
				os.Exit(-1)
			}
			msg := messageColors[0].Sprint(fmt.Sprintf(messageFormat, i))
			tracker.UpdateMessage(msg)
			i++
		}
	}
}
