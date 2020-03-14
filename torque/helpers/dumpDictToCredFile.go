package helpers

import (
	"fmt"
	"io/ioutil"
	"torque/customTypes"
)

type CredDict = customTypes.CredDict

func DumpDictToCredFile(fileLocation string, dictData map[string]CredDict) {
	data := ""
	counter := 0
	for profile, _ := range dictData {
		if counter == 0 {
			data = data + "[" + string(profile) + "]"
			counter = counter + 1
		} else {
			data = data + "\n[" + string(profile) + "]"
		}
		keys := CredDict{}
		keys = dictData[profile]
		data = data + "\naws_access_key_id = " + keys.AccessKey
		data = data + "\naws_secret_access_key = " + keys.SecretKey
		if keys.SessionToken != "" {
			data = data + "\naws_session_token = " + keys.SessionToken
		}
	}
	d1 := []byte(data)
	error := ioutil.WriteFile(fileLocation, d1, 0644)
	if error != nil {
		fmt.Println("[-] Error occurred while writing to file")
		fmt.Println(error)
	} else {
		fmt.Println("[+] Successfully written to file : " + fileLocation)
	}
}
