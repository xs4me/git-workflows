package utils

import (
	"gepaplexx/git-workflows/logger"
)

const AlreadyUpToDateErr = "already up-to-date"

func CheckIfError(err error) {
	if err != nil {
		if err.Error() == AlreadyUpToDateErr {
			logger.Info(AlreadyUpToDateErr)
			return
		}
		logger.Fatal(err.Error())
		panic(err)
	}
}
