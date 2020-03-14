package authMFA

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"os"
	"os/user"
	"strings"
	"torque/customTypes"
	"torque/helpers"
)

type CredDict = customTypes.CredDict

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
	credFileData = helpers.ReadCredsFile(cwd)
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
	helpers.DumpDictToCredFile(cwd, credFileData)
}
