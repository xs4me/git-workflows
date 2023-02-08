package api

import (
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"os"
)

func gitAuthMethod(c *model.Config) transport.AuthMethod {
	privateKeyfile := c.SshConfigDir + "id_rsa"
	_, err := os.Stat(privateKeyfile)
	if err != nil {
		logger.Fatal("SSH key not found. Please provide a valid SSH key. Password authentication is not yet supported.")
		os.Exit(1)
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyfile, "")
	utils.CheckIfError(err)

	return publicKeys
}
