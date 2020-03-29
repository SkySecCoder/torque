package authMFA

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sts"
	"io/ioutil"
	"os"
	"os/user"
	"torque/helpers"
)

func dumpCache(profile string, credentials *sts.GetSessionTokenOutput) {
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

	rawData, _ := json.MarshalIndent(credentials, "", "    ")

	err = ioutil.WriteFile(cachePath+profile+"-credentials.json", []byte(rawData), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
