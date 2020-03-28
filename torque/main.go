package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"torque/authMFA"
	"torque/customTypes"
	"torque/helpers"
	"torque/keyRotation"
	"torque/programHelp"
)

type CredDict = customTypes.CredDict

func main() {
	progArgs := os.Args
	firstArgument := strings.Split(progArgs[0], "/")
	programName := firstArgument[len(firstArgument)-1]
	test := true
	if test == false {
	if len(progArgs) == 1 {
		programHelp.ProgramHelp(programName)
	} else if progArgs[1] == "help" {
		programHelp.ProgramHelp(programName)
	} else if progArgs[1] == "rotate" {
		cwd := helpers.GetAWSCredentialFileLocation()
		exists, error := helpers.DoesFileExist(cwd)
		if error == nil {
			if exists == true {
				// Checking program args
				if progArgs[2] == "all" {
					rotateAll()
				} else if len(progArgs) == 3 {
					if strings.Contains(progArgs[2], "mfa-") {
						fmt.Println("[-] Not rotating profile " + progArgs[2])
						fmt.Println("[-] Cannot rotate profiles that contain 'mfa-' in their name...")
					} else {
						keyRotation.RotateKey(progArgs[2])
					}
				} else {
					programHelp.ProgramHelp(programName)
				}
			} else {
				fmt.Println("[-] Unable to find creds file at : " + cwd)
			}
		} else {
			fmt.Println("[-] Error occurred while finding creds file at : " + cwd)
		}
	} }


	if len(os.Args) == 1 || os.Args[1] == "help" || os.Args[1] == "-h" {
		programHelp.ProgramHelp(programName)
	} else if os.Args[1] == "auth" && len(os.Args) == 3 {
		authMFA.AuthMFA(progArgs[2])
	} else if os.Args[1] == "rotate" && len(os.Args) == 3 {
		fmt.Println("Oops...")
	} else {
		programHelp.ProgramHelp(programName)
	}
}

func rotateAll() {
	cwd := helpers.GetAWSCredentialFileLocation()
	credFileData := helpers.ReadCredsFile(cwd)
	for profile, _ := range credFileData {
		if credFileData[profile].SessionToken == "" {
			keyRotation.RotateKey(profile)
		} else {
			fmt.Println("\n[-] Not rotating " + profile + "\n")
		}
	}
}

func printJsonArray(data []string) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("[-] Error occurred file converting to json")
	} else {
		fmt.Println(string(jsonData))
	}
}

func printJsonCustom(data map[string]CredDict) {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("[-] Error occurred file converting to json")
	} else {
		fmt.Println(string(jsonData))
	}
}
