package programHelp

import (
	"fmt"
)

func ProgramHelp(programName string) {
	data := "\nUsage: " + programName + " [OPTION] [ARGUMENT]"
	data = data + "\nUsed to manage AWS access keys on local system.\n"
	data = data + "\n\thelp,\t\t\t\tshows program help"
	data = data + "\n\trotate [PROFILE_NAME],\t\trotates access keys for [PROFILE_NAME]"
	data = data + "\n\t\t\tall,\t\trotates all access keys in '$HOME/.aws/credentials' file"
	data = data + "\n\tauth [PROFILE_NAME],\t\tauths mfa for [PROFILE_NAME]"
	data = data + "\n\tassume [PROFILE_NAME],\t\tassumes IAM role as specified in [PROFILE_NAME]"
	data = data + "\n"
	fmt.Println(data)
}
