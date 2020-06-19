package assumeRole

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"io/ioutil"
	"os"
	"strings"
	"torque/helpers"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-ini/ini"
	"os/user"
	"torque/authMFA"
)

func AssumeRole(profile string) {
	profile = strings.ReplaceAll(profile, "\n", "")
	config := readAWSConfig(profile)
	sourceProfile := config["source_profile"]

	// Checking if the creds, for the the role that we want to assume, have expired
	credsNotExpired, cacheCreds := helpers.CheckCache(profile)
	if credsNotExpired {
		fmt.Println("export AWS_ACCESS_KEY_ID=" + cacheCreds.AccessKey)
		fmt.Println("export AWS_SECRET_ACCESS_KEY=" + cacheCreds.SecretKey)
		fmt.Println("export AWS_SESSION_TOKEN=" + cacheCreds.SessionToken)
		return
	}

	// Checking if the creds, for the the source profile, have expired
	credsNotExpired, cacheCreds = helpers.CheckCache(sourceProfile)
	if credsNotExpired == false {
		authMFA.AuthMFA(sourceProfile, "silent")
		_, cacheCreds = helpers.CheckCache(sourceProfile)
	}

	mysession, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(cacheCreds.AccessKey, cacheCreds.SecretKey, cacheCreds.SessionToken),
	})

	// Getting caller identity
	stsClient := sts.New(mysession)
	result, err := stsClient.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(config["role_arn"]),
		RoleSessionName: aws.String("aakash"),
	})
	if err != nil {
		fmt.Println(err)
	}

	dumpCache(profile, result)
	fmt.Println("export AWS_ACCESS_KEY_ID=" + string(*result.Credentials.AccessKeyId))
	fmt.Println("export AWS_SECRET_ACCESS_KEY=" + string(*result.Credentials.SecretAccessKey))
	fmt.Println("export AWS_SESSION_TOKEN=" + string(*result.Credentials.SessionToken))
}

func readAWSConfig(profile string) map[string]string {
	iniReader, err := ini.Load(helpers.GetAWSConfigFileLocation())
	if err != nil {
		fmt.Println(err)
	}
	fileData, err := iniReader.GetSection("profile " + profile)
	if err != nil {
		fmt.Println(err)
	}
	returnData := map[string]string{
		"mfa_serial":     "",
		"role_arn":       "",
		"source_profile": "",
	}
	returnData["mfa_serial"] = fileData.Key("mfa_serial").String()
	returnData["role_arn"] = fileData.Key("role_arn").String()
	returnData["source_profile"] = fileData.Key("source_profile").String()
	return returnData
}

func dumpCache(profile string, credentials *sts.AssumeRoleOutput) {
	homedir := ""
	user, err := user.Current()
	if err == nil {
		homedir = user.HomeDir
	} else {
		fmt.Println(err)
		return
	}

	cachePath := homedir + "/.torque/cache/"
	if exist, _ := helpers.DoesFileExist(cachePath); exist == false {
		err = os.MkdirAll(cachePath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	rawData, _ := json.MarshalIndent(map[string]*sts.Credentials{"Credentials": credentials.Credentials}, "", "    ")

	err = ioutil.WriteFile(cachePath+profile+"-credentials.json", []byte(rawData), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
