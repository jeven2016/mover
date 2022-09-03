package common

import "github.com/jedib0t/go-pretty/v6/progress"

type Parameters struct {
	From                string `mapstructure:"from"`
	To                  string `mapstructure:"to"`
	FileExtension       string `mapstructure:"file_extension"`
	MinSize             string `mapstructure:"file_min_size"`
	CheckPicture        bool   `mapstructure:"check_picture"`
	CreateRootDirectory bool   `mapstructure:"create_root_directory"`
}

type Stats struct {
	FileCount    int
	IgnoreFiles  int
	PictureCount int
	FolderCount  int
}

type ProgressSetting struct {
	Total int64
	Units *progress.Units
}
