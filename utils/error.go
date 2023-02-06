package utils

func CheckIfError(err error) {
	if err != nil {
		panic(err)
	}
}
