package keyRotation

import (
	"fmt"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os/user"
	"io/ioutil"
)

type CredDict struct {
	AccessKey		string 	`json:"accessKey"`
	SecretKey 		string 	`json:"secretKey"`
	SessionToken 	string 	`json:"sessionToken"`
}

func RotateKey(profile string) {
	fmt.Println("\n[+] Rotating credentials for profile : "+profile+"\n")
	credFileData := map[string]CredDict{}
	// Getting credential file location
	cwd := ""
	user, err := user.Current()
	if err == nil {
		cwd = user.HomeDir+"/.aws/credentials"
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

	// Setting Credentials
	//Delete this if operationFailure == false {
	creds := credentials.NewSharedCredentials(cwd, profile)
	credValue, err := creds.Get()
	if err != nil {
		fmt.Println("[-] Cannot load creds for profile : "+profile)
		fmt.Println()
		fmt.Println(err)
		fmt.Println()
		return
	} else {
		fmt.Println("[+] Successfully loaded creds")
	}

	// Delete This Key
	keyToDelete := ""
	if strings.Contains(profile, "mfa-") {
		myKey := credFileData[strings.ReplaceAll(profile, "mfa-", "")]
		keyToDelete = myKey.AccessKey
	} else {
		keyToDelete = credValue.AccessKeyID
	}
		

	// Creating Session
	mysession, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
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

	// Creating new access key
	newKey := CredDict{}

	// Creating IAM client
	iamClient := iam.New(mysession)

	// Creating new access key
	resultIam, err := iamClient.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(myuser),
	})
	if err != nil {
		fmt.Println("[-] Failed in creating new access key")
		fmt.Println()
		if strings.Contains(err.Error(), "LimitExceeded: Cannot exceed quota for AccessKeysPerUser: 2") {
			fmt.Println("[-] It appears you have 2 access keys...")
			fmt.Println("[-] This program will not work for 2 access keys...")
			return
		} else {
			fmt.Println(err)
			return
		}
		fmt.Println()
		return
	} else {
		fmt.Println("[+] Successfully created new access key")
		//fmt.Println(*resultIam.AccessKey)
		keyDetails := *resultIam.AccessKey
		newKey.AccessKey = *keyDetails.AccessKeyId
		newKey.SecretKey = *keyDetails.SecretAccessKey
		newKey.SessionToken = ""
	}

	fmt.Println("[-] Deleting old access key : "+keyToDelete)

	// Deleting previous key
	_, err = iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
		UserName: aws.String(myuser),
		AccessKeyId: aws.String(keyToDelete),
	})
	if err != nil {
		fmt.Println("[-] Failed to delete : "+keyToDelete)
		fmt.Println()
		fmt.Println(err)
		fmt.Println()
		return
	} else {
		fmt.Println("[+] Successfully deleted : "+keyToDelete)
	}
	credFileData[profile] = newKey
	dumpDictToCredFile(cwd, credFileData)
}

func readCredsFile(cwd string) map[string]CredDict{
	returnData := map[string]CredDict{}

	// Getting credential file location
	
	data, err := ioutil.ReadFile(cwd)
	if err != nil {
		fmt.Println("[-] Error occurred while reading the file "+cwd)
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
			data = data+"["+string(profile)+"]"
			counter = counter + 1
		} else {
			data = data+"\n["+string(profile)+"]"
		}
		keys := CredDict{}
		keys = dictData[profile]
		data = data+"\naws_access_key_id = "+keys.AccessKey
		data = data+"\naws_secret_access_key = "+keys.SecretKey
		if keys.SessionToken != ""{
			data = data+"\naws_session_token = "+keys.SessionToken
		}
	}
	d1 := []byte(data)
	error := ioutil.WriteFile(fileLocation, d1, 0644)
	if error != nil {
		fmt.Println("[-] Error occurred while writing to file")
		fmt.Println(error)
	} else {
		fmt.Println("[+] Successfully written to file : "+fileLocation)
	}
}

func convertArrayToMap(data []string) map[string]CredDict {
	returnData := map[string]CredDict{}
	profile := ""
	keyData := CredDict{}

	for line := range data {
		if (strings.Contains(data[line], "[") && strings.Contains(data[line], "]")) {
			keyData.AccessKey = ""
			keyData.SecretKey = ""
			keyData.SessionToken = ""

			profile = data[line]
			profile = strings.ReplaceAll(profile, "[", "")
			profile = strings.ReplaceAll(profile, "]", "")
		} else if (strings.Contains(data[line], "aws_access_key_id = ")) {
			accessKey := data[line]
			accessKey = strings.ReplaceAll(accessKey, "aws_access_key_id = ", "")
			keyData.AccessKey = accessKey
		} else if (strings.Contains(data[line], "aws_secret_access_key = ")) {
			secretKey := data[line]
			secretKey = strings.ReplaceAll(secretKey, "aws_secret_access_key = ", "")
			keyData.SecretKey = secretKey
			returnData[profile] = keyData
		} else if (strings.Contains(data[line], "aws_session_token = ")) {
			sessionToken := data[line]
			sessionToken = strings.ReplaceAll(sessionToken, "aws_session_token = ", "")
			keyData.SessionToken = sessionToken
			returnData[profile] = keyData
		}
	}
	return returnData
}
