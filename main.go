/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import "gepaplexx/git-workflows/cmd"

// version will be set from the build command and overwrites the default value
var Version string = "0.0.1"

func main() {
	cmd.Version = Version
	cmd.Execute()
}
