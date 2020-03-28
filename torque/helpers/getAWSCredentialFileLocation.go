package helpers

import (
	"os/user"
	"fmt"
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
