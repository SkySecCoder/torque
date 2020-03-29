package helpers

import (
	"io/ioutil"
	"fmt"
	"os/user"
	"encoding/json"
	"time"
)

type CacheData struct {
	AccessKeyId		string
	Expiration		string
	SecretAccessKey	string
	SessionToken	string
}

func CheckCache(profile string) (bool, CredDict) {
	homedir := ""
	user, err := user.Current()
	if err == nil {
		homedir = user.HomeDir
	} else {
		fmt.Println(err)
	}

	rawData,_ := ioutil.ReadFile(homedir+"/.torque/cache/"+profile+"-credentials.json")

	var data map[string]CacheData
	_ = json.Unmarshal(rawData, &data)

	currentTime := time.Now()
	expTime,_ := time.Parse("2006-01-02T15:04:05Z", data["Credentials"].Expiration)

	returnKey := CredDict{}
	returnKey.AccessKey = data["Credentials"].AccessKeyId
	returnKey.SecretKey = data["Credentials"].SecretAccessKey
	returnKey.SessionToken = data["Credentials"].SessionToken
	return currentTime.Before(expTime), returnKey
}
