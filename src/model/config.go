package model

import "fmt"

type Config struct {
	Development    bool   `json:"development"`
	LegacyBehavior bool   `json:"legacy_behavior"` // if true falls back to multibranch behavior
	BaseDir        string `json:"base_dir"`
	GitUrl         string `json:"git_url"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Reponame       string `json:"reponame"`
	Branch         string `json:"branch"`
	Extract        bool   `json:"extract"`
}

func (c *Config) ClonePath() string {

	return fmt.Sprintf("%s/%s", c.BaseDir, c.Reponame)
}
