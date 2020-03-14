package helpers

import (
	"io/ioutil"
	"strings"
	"fmt"
)

func ReadCredsFile(cwd string) map[string]CredDict {
	returnData := map[string]CredDict{}

	// Getting credential file location

	data, err := ioutil.ReadFile(cwd)
	if err != nil {
		fmt.Println("[-] Error occurred while reading the file " + cwd)
	} else {
		rawCreds := strings.Split(string(data), "\n")
		returnData = ConvertArrayToMap(rawCreds)
	}
	return returnData
}
