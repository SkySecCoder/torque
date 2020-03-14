package keyRotation

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"os"
	"os/user"
	"strings"
	"torque/authMFA"
	"torque/customTypes"
	"torque/helpers"
)

type CredDict = customTypes.CredDict

func RotateKey(profile string) {
	fmt.Println("\n[+] Rotating credentials for profile : " + profile + "\n")
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
	credFileData = helpers.ReadCredsFile(cwd)
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
		fmt.Println("[-] Cannot load creds for profile : " + profile)
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
		if strings.Contains(err.Error(), "explicit deny\n	status code: 403") {
			fmt.Print("[-] It appears you have an explicit deny for this operations\n[?] Does " + profile + " require MFA authentication(y/n) : ")
			reader := bufio.NewReader(os.Stdin)
			option, _ := reader.ReadString('\n')
			if strings.ReplaceAll(option, "\n", "") == "y" {
				rotateWithMFA(profile, cwd)
				return
			} else {
				fmt.Println("[-] Exiting...")
				return
			}
		} else if strings.Contains(err.Error(), "LimitExceeded: Cannot exceed quota for AccessKeysPerUser: 2") {
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

	fmt.Println("[-] Deleting old access key : " + keyToDelete)

	// Deleting previous key
	_, err = iamClient.DeleteAccessKey(&iam.DeleteAccessKeyInput{
		UserName:    aws.String(myuser),
		AccessKeyId: aws.String(keyToDelete),
	})
	if err != nil {
		fmt.Println("[-] Failed to delete : " + keyToDelete)
		fmt.Println()
		fmt.Println(err)
		fmt.Println()
		return
	} else {
		fmt.Println("[+] Successfully deleted : " + keyToDelete)
	}
	credFileData[profile] = newKey
	helpers.DumpDictToCredFile(cwd, credFileData)
}



func rotateWithMFA(profile string, cwd string) {
	authMFA.AuthMFA(profile)
	RotateKey("mfa-" + profile)
	credData := helpers.ReadCredsFile(cwd)
	delete(credData, profile)
	credData[profile] = credData["mfa-"+profile]

	key := credData[profile]
	key.SessionToken = ""
	credData[profile] = key

	fmt.Println("[+] Deleting profile : mfa-" + profile)
	delete(credData, "mfa-"+profile)
	fmt.Println("\n[+] Successfully rotated MFA creds for : " + profile + "\n")
	helpers.DumpDictToCredFile(cwd, credData)
}


