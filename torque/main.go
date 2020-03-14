package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"
	"torque/authMFA"
	"torque/customTypes"
	"torque/keyRotation"
	"torque/helpers"
)

type CredDict = customTypes.CredDict

func main() {
	progArgs := os.Args
	firstArgument := strings.Split(progArgs[0], "/")
	programName := firstArgument[len(firstArgument)-1]
	if len(progArgs) == 1 {
		programHelp(programName)
	} else if progArgs[1] == "help" {
		programHelp(programName)
	} else if progArgs[1] == "rotate" {
		cwd := getAWSCredentialFileLocation()
		exists, error := doesFileExist(cwd)
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
					programHelp(programName)
				}
			} else {
				fmt.Println("[-] Unable to find creds file at : " + cwd)
			}
		} else {
			fmt.Println("[-] Error occurred while finding creds file at : " + cwd)
		}
	} else if progArgs[1] == "auth" {
		cwd := getAWSCredentialFileLocation()
		exists, error := doesFileExist(cwd)
		if error == nil {
			if exists == true {
				// Checking program args
				if len(progArgs) == 3 {
					authMFA.AuthMFA(progArgs[2])
				} else {
					programHelp(programName)
				}
			} else {
				fmt.Println("[-] Unable to find creds file at : " + cwd)
			}
		} else {
			fmt.Println("[-] Error occurred while finding creds file at : " + cwd)
		}
	} else {
		fmt.Println("[-] Wrong option, exiting...")
	}
}

func programHelp(programName string) {
	data := "\nUsage: " + programName + " [OPTION] [ARGUMENT]"
	data = data + "\nUsed to manage AWS access keys on local system.\n"
	data = data + "\n\thelp,\t\t\t\tshows program help"
	data = data + "\n\trotate [PROFILE_NAME],\t\trotates access keys for [PROFILE_NAME]"
	data = data + "\n\t\t\tall,\t\trotates all access keys in '$HOME/.aws/credentials' file"
	data = data + "\n\tauth [PROFILE_NAME],\t\tauths mfa for [PROFILE_NAME]"
	data = data + "\n"
	fmt.Println(data)
}

func rotateAll() {
	cwd := getAWSCredentialFileLocation()
	credFileData := helpers.ReadCredsFile(cwd)
	for profile, _ := range credFileData {
		if credFileData[profile].SessionToken == "" {
			keyRotation.RotateKey(profile)
		} else {
			fmt.Println("\n[-] Not rotating " + profile + "\n")
		}
	}
}

func getAWSCredentialFileLocation() string {
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

func doesFileExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
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
