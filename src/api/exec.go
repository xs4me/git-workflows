package api

import (
	"bytes"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/utils"
	"os/exec"
	"strings"
)

func execute(cmd *exec.Cmd) string {
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	logger.Debug(cmd.String())
	err := cmd.Run()
	utils.CheckIfError(err)

	return strings.TrimRight(out.String(), "\n")
}
