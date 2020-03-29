package authMFA

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"os"
	"strings"
	"torque/customTypes"
	"torque/helpers"
)

type CredDict = customTypes.CredDict

func AuthMFA(profile string, mode string) {
	credsNotExpired, cacheCreds := helpers.CheckCache(profile)
	if credsNotExpired {
		if mode != "silent" {
			fmt.Println("export AWS_ACCESS_KEY_ID=" + cacheCreds.AccessKey)
			fmt.Println("export AWS_SECRET_ACCESS_KEY=" + cacheCreds.SecretKey)
			fmt.Println("export AWS_SESSION_TOKEN=" + cacheCreds.SessionToken)
		}
		return
	}
	cwd := helpers.GetAWSCredentialFileLocation()
	exists, error := helpers.DoesFileExist(cwd)

	if error != nil {
		fmt.Println(error)
		return
	}

	if exists != true {
		fmt.Println("[-] Unable to find creds file at : " + cwd)
		return
	}

	// Asking user for their MFA Token
	fmt.Print("[+] Please enter you MFA token code : ")
	reader := bufio.NewReader(os.Stdin)
	mfaToken, _ := reader.ReadString('\n')
	//fmt.Println(mfaToken)

	credFileData := map[string]CredDict{}

	// Checking if profile even exists
	credFileData = helpers.ReadCredsFile(cwd)
	_, ok := credFileData[profile]
	if ok == false {
		fmt.Println("[-] This profile does not exist in the creds file")
		return
	}

	creds := credentials.NewSharedCredentials(cwd, profile)
	_, err := creds.Get()
	if err != nil {
		fmt.Println("[-] Cannot load creds for profile : " + profile)
		fmt.Println()
		fmt.Println(err)
		fmt.Println()
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

	// Getting token using MFA
	response, err := stsClient.GetSessionToken(&sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(43200),
		SerialNumber:    aws.String(requiredArn + "mfa/" + myuser),
		TokenCode:       aws.String(strings.ReplaceAll(mfaToken, "\n", "")),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	dumpCache(profile, response)

	if mode != "silent" {
		fmt.Println("export AWS_ACCESS_KEY_ID=" + *response.Credentials.AccessKeyId)
		fmt.Println("export AWS_SECRET_ACCESS_KEY=" + *response.Credentials.SecretAccessKey)
		fmt.Println("export AWS_SESSION_TOKEN=" + *response.Credentials.SessionToken)
	}
}
