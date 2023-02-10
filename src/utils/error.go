package utils

import (
	"gepaplexx/git-workflows/logger"
)

func CheckIfError(err error) {
	if err != nil {
		logger.Fatal(err.Error())
		panic(err)
	}
}

func CheckCommandError(err error, out string) {
	if err != nil {
		logger.Fatal(out)
		panic(err)
	}
}
