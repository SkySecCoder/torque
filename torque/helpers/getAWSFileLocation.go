package helpers

import (
	"fmt"
	"os/user"
)

func GetAWSCredentialFileLocation() string {
	cwd := ""
	user, err := user.Current()
	if err == nil {
		cwd = user.HomeDir
	} else {
		fmt.Println(err)
	}
	cwd = cwd + "/.aws/credentials"

	return (cwd)
}

func GetAWSConfigFileLocation() string {
	cwd := ""
	user, err := user.Current()
	if err == nil {
		cwd = user.HomeDir
	} else {
		fmt.Println(err)
	}
	cwd = cwd + "/.aws/config"

	return (cwd)
}
