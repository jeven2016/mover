package common

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"math/rand"
	"os"
	"time"
)

func ShowProgress() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Task", "Result"})
	//----------------------

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

	messageColors := []text.Color{
		text.FgRed,
		text.FgGreen,
		text.FgYellow,
		text.FgBlue,
		text.FgMagenta,
		text.FgCyan,
		text.FgWhite,
	}

	go pw.Render()

	messageFormat := "Downloading Files: %3d"
	tracker := &progress.Tracker{Message: "Progress text", Total: 50, Units: progress.UnitsDefault}

	pw.AppendTracker(tracker)

	tick := time.Tick(300 * time.Millisecond)
	i := 1
	for !tracker.IsDone() {
		select {

		case <-tick:
			tracker.Increment(1)
			if i >= 50 {
				print("All Done")
				renderTable(t, i, rand.Intn(20))
				tracker.MarkAsDone()
				os.Exit(-1)
			}
			msg := messageColors[0].Sprint(fmt.Sprintf(messageFormat, i))
			tracker.UpdateMessage(msg)
			i++
		}
	}
}

func renderTable(t table.Writer, success int, failue int) {
	t.ResetRows()
	t.AppendRows([]table.Row{
		{1, "File Downloaded", success},
	})
	t.AppendFooter(table.Row{2, "Failure", failue})
	t.Render()
}
