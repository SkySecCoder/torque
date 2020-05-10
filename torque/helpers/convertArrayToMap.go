package helpers

import (
	"strings"
)

func ConvertArrayToMap(data []string) map[string]CredDict {
	returnData := map[string]CredDict{}
	profile := ""
	keyData := CredDict{}

	for line := range data {
		data[line] = strings.ReplaceAll(data[line], " ", "")
		if strings.Contains(data[line], "[") && strings.Contains(data[line], "]") {
			keyData.AccessKey = ""
			keyData.SecretKey = ""
			keyData.SessionToken = ""

			profile = data[line]
			profile = strings.ReplaceAll(profile, "[", "")
			profile = strings.ReplaceAll(profile, "]", "")
		} else if strings.Contains(data[line], "aws_access_key_id=") {
			accessKey := data[line]
			accessKey = strings.ReplaceAll(accessKey, "aws_access_key_id=", "")
			keyData.AccessKey = accessKey
		} else if strings.Contains(data[line], "aws_secret_access_key=") {
			secretKey := data[line]
			secretKey = strings.ReplaceAll(secretKey, "aws_secret_access_key=", "")
			keyData.SecretKey = secretKey
			returnData[profile] = keyData
		} else if strings.Contains(data[line], "aws_session_token=") {
			sessionToken := data[line]
			sessionToken = strings.ReplaceAll(sessionToken, "aws_session_token=", "")
			keyData.SessionToken = sessionToken
			returnData[profile] = keyData
		}
	}
	return returnData
}
