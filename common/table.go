package common

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

func ShowTable(stats *Stats) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Stats", "Count"})
	t.AppendRows([]table.Row{
		{"Files Moved", stats.FileCount},
		{"Files Ignored", stats.IgnoreFiles},
		{"Pictures Moved", stats.PictureCount},
		{"All Folders", stats.FolderCount},
		{"Picture Moved", stats.PictureCount},
		{"Picture Moved Failure", stats.PictureFailure},
		{"Folders Removed", stats.RemoveOldFolders},
		{"Picture Removed Failure", stats.RemoveOldFoldersFailure},
	})
	//t.AppendFooter(table.Row{2, "Failure", failue})
	t.Render()
}
