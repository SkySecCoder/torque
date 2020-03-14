package authMFA

import (
	"fmt"
	"bufio"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"strings"
	"os/user"
	"io/ioutil"
)

type CredDict struct {
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	SessionToken string `json:"sessionToken"`
}

func AuthMFA(profile string) {
	// Asking user for their MFA Token
	fmt.Print("\n[+] Please enter you MFA token code : ")
	reader := bufio.NewReader(os.Stdin)
	mfaToken, _ := reader.ReadString('\n')
	//fmt.Println(mfaToken)

	credFileData := map[string]CredDict{}
	// Getting credential file location
	cwd := ""
	user, err := user.Current()
	if err == nil {
		cwd = user.HomeDir + "/.aws/credentials"
	} else {
		fmt.Println(err)
	}

	// Checking if profile even exists
	credFileData = readCredsFile(cwd)
	_, ok := credFileData[profile]
	if ok == false {
		fmt.Println("[-] This profile does not exist in the creds file")
		return
	}

	creds := credentials.NewSharedCredentials(cwd, profile)
	_, err = creds.Get()
	if err != nil {
		fmt.Println("[-] Cannot load creds for profile : " + profile)
		fmt.Println()
		fmt.Println(err)
		fmt.Println()
	} else {
		fmt.Println("[+] Successfully loaded creds")
	}

	// Creating Session
	mysession, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	})

	// Getting caller identity
	stsClient := sts.New(mysession)
	result, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Println(err)
	}

	arn := *result.Arn
	myuser := arn[31:]
	requiredArn := arn[:26]
	//fmt.Println(arn)
	//fmt.Println(myuser)
	//fmt.Println(arn[:26])

	// Getting token using MFA
	response, err := stsClient.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(3600),
		SerialNumber:    aws.String(requiredArn + "mfa/" + myuser),
		TokenCode:       aws.String(strings.ReplaceAll(mfaToken, "\n", "")),
	})
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("[+] Successfully authenticated MFA")
	}
	newKey := CredDict{}
	newKey.AccessKey = *response.Credentials.AccessKeyId
	newKey.SecretKey = *response.Credentials.SecretAccessKey
	newKey.SessionToken = *response.Credentials.SessionToken
	credFileData["mfa-"+profile] = newKey
	dumpDictToCredFile(cwd, credFileData)
}

func readCredsFile(cwd string) map[string]CredDict {
	returnData := map[string]CredDict{}

	// Getting credential file location

	data, err := ioutil.ReadFile(cwd)
	if err != nil {
		fmt.Println("[-] Error occurred while reading the file " + cwd)
	} else {
		rawCreds := strings.Split(string(data), "\n")
		returnData = convertArrayToMap(rawCreds)
	}
	return returnData
}

func dumpDictToCredFile(fileLocation string, dictData map[string]CredDict) {
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

func convertArrayToMap(data []string) map[string]CredDict {
	returnData := map[string]CredDict{}
	profile := ""
	keyData := CredDict{}

	for line := range data {
		if strings.Contains(data[line], "[") && strings.Contains(data[line], "]") {
			keyData.AccessKey = ""
			keyData.SecretKey = ""
			keyData.SessionToken = ""

			profile = data[line]
			profile = strings.ReplaceAll(profile, "[", "")
			profile = strings.ReplaceAll(profile, "]", "")
		} else if strings.Contains(data[line], "aws_access_key_id = ") {
			accessKey := data[line]
			accessKey = strings.ReplaceAll(accessKey, "aws_access_key_id = ", "")
			keyData.AccessKey = accessKey
		} else if strings.Contains(data[line], "aws_secret_access_key = ") {
			secretKey := data[line]
			secretKey = strings.ReplaceAll(secretKey, "aws_secret_access_key = ", "")
			keyData.SecretKey = secretKey
			returnData[profile] = keyData
		} else if strings.Contains(data[line], "aws_session_token = ") {
			sessionToken := data[line]
			sessionToken = strings.ReplaceAll(sessionToken, "aws_session_token = ", "")
			keyData.SessionToken = sessionToken
			returnData[profile] = keyData
		}
	}
	return returnData
}

