package model

type Config struct {
	LegacyBehavior bool   `json:"legacy_behavior"` // if true falls back to multibranch behavior
	BaseDir        string `json:"base_dir"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	//InfraRepoSuffix  string `json:"infrarepo_suffix"`
	//ImageTagLocation string `json:"image_tag_location"`
	//
	//Branch             string `json:"branch"`
	//Environment        string `json:"environment"`
	//Repository         string `json:"repository"`
	//GitUrl             string `json:"git_url"`
	//Namespace          string `json:"namespace"`
	//DeploySourceBranch string `json:"deploy_source_branch"`
	//DeployTargetBranch string `json:"deploy_target_branch"`
}
