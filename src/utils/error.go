package utils

import (
	"gepaplexx/git-workflows/logger"
)

const AlreadyUpToDateErr = "already up-to-date"
const ElementNotFoundErr = "element not found"

func CheckIfError(err error) {
	if err != nil {
		if err.Error() == AlreadyUpToDateErr || err.Error() == ElementNotFoundErr {
			logger.Info(err.Error())
			return
		}
		logger.Fatal(err.Error())
		panic(err)
	}
}
