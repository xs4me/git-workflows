package model

import "fmt"

type Config struct {
	Development bool `json:"development"`

	LegacyBehavior bool   `json:"legacy_behavior"` // if true falls back to multibranch behavior
	BaseDir        string `json:"base_dir"`
	GitUrl         string `json:"git_url"`
	Username       string `json:"username"`
	Email          string `json:"email"`

	Reponame string `json:"reponame"`
	//InfraRepoSuffix  string `json:"infrarepo_suffix"`
	//ImageTagLocation string `json:"image_tag_location"`
	//
	//Branch             string `json:"branch"`
	//Environment        string `json:"environment"`
	//Repository         string `json:"repository"`

	//Namespace          string `json:"namespace"`
	//DeploySourceBranch string `json:"deploy_source_branch"`
	//DeployTargetBranch string `json:"deploy_target_branch"`
}

func (c *Config) LocalPath() string {
	return fmt.Sprintf("%s/%s", c.BaseDir, c.Reponame)
}
