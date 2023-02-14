package model

import (
	"fmt"
)

type Config struct {
	Development     bool     `json:"development"`
	BaseDir         string   `json:"base_dir"`
	GitUrl          string   `json:"git_url"`
	Username        string   `json:"username"`
	Email           string   `json:"email"`
	Reponame        string   `json:"reponame"`
	Branch          string   `json:"branch"`
	Extract         bool     `json:"extract"`
	SshConfigDir    string   `json:"ssh_dir" `
	RepoToken       string   `json:"repo_token"`
	InfraRepoSuffix string   `json:"infra_repo_suffix"`
	ImageTag        string   `json:"image_tag"`
	TagLocation     string   `json:"tag_location"`
	Stages          []string `json:"stages"`
	Env             string   `json:"env"`
	FromBranch      string   `json:"from_branch"`
}

func (c *Config) ApplicationClonePath() string {

	return fmt.Sprintf("%s/%s", c.BaseDir, c.Reponame)
}

func (c *Config) InfrastructureClonePath() string {
	return fmt.Sprintf("%s/%s%s", c.BaseDir, c.Reponame, c.InfraRepoSuffix)
}

func (c *Config) ImageTagLocation() string {
	if "" == c.TagLocation {
		return fmt.Sprintf(".%s.image.tag", c.Reponame)
	} else {
		return c.TagLocation
	}
}

func (c *Config) IsPushEnabled() bool {
	return !c.Development
}
