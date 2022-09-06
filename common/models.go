package common

import "github.com/jedib0t/go-pretty/v6/progress"

type GlobalSetting struct {
	Parameters *Parameters `mapstructure:"global_setting"`
}

type Parameters struct {
	From                string `mapstructure:"from"`
	To                  string `mapstructure:"to"`
	FileExtension       string `mapstructure:"file_extension"`
	MinSize             string `mapstructure:"file_min_size"`
	PicExtension        string `mapstructure:"picture_extension"`
	PicMinSize          string `mapstructure:"pic_min_size"`
	CheckPicture        bool   `mapstructure:"check_picture"`
	CreateRootDirectory bool   `mapstructure:"create_root_directory"`
}

type Stats struct {
	FileCount               int32
	IgnoreFiles             int32
	FileFailure             int32
	PictureCount            int32
	PictureFailure          int32
	FolderCount             int32
	RemoveOldFolders        int32
	RemoveOldFoldersFailure int32
}

type ProgressSetting struct {
	Total int64
	Units *progress.Units
}

type ResourceType int32

const (
	FileType ResourceType = iota
	PicType
)
