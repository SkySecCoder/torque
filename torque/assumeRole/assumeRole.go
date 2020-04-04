package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"torque/helpers"
	//"io/ioutil"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go/service/stscreds"
	"github.com/go-ini/ini"
)

func main() {
	fmt.Printf("Profile you want to assume : ")
	input := bufio.NewReader(os.Stdin)
	profile, err := input.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	profile = strings.ReplaceAll(profile, "\n", "")
	readAWSConfig(profile)
	/*config := readAWSConfig(profile)
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "myProfile",
	})*/
}

func readAWSConfig(profile string) {
	iniReader,err := ini.Load(helpers.GetAWSConfigFileLocation())
	if err != nil {
		fmt.Println(err)
	}
	fileData,err := iniReader.GetSection("profile "+profile)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fileData.Key("mfa_serial"))
}
