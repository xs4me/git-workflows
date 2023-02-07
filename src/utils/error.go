package utils

import (
	"gepaplexx/git-workflows/logger"
)

func CheckIfError(err error) {
	if err != nil {
		logger.Fatal(err.Error())
	}
}
