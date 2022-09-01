package common

type Parameters struct {
	From                string `mapstructure:"from"`
	To                  string `mapstructure:"to"`
	FileExtension       string `mapstructure:"file_extension"`
	MinSize             string `mapstructure:"file_min_size"`
	CheckPicture        bool   `mapstructure:"check_picture"`
	CreateRootDirectory bool   `mapstructure:"create_root_directory"`
}
