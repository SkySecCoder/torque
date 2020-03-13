# Torque
### Introduction ###
Torque is a tool for rotating AWS access keys as well as using the CLI in order to use those access keys with MFA.
### Objective ###
AWS by deafult allows the use of it's access keys without the user needing to use their MFA.

This is a security flaw especially if the keys are leaked to unauthorised parties.

These can be leaked due to accidental code push to GitHub, or simply by copying the unencrypted keys from the `~/.aws/credentials` file to an unauthorised location.

In order to protect these keys, the IAM user must have a policy that forces them to use MFA while using their keys. This can be achieved by enfrocing MFA via an IAM policy and using `Torque` to create temporary access keys that the user can use for their work.

Unfortunately Torque cannot be used when dealing with service users who may need the key but are enable to provide their MFA.(Note: you should make sure that the service user keys are regularly rotated, are not accessible by many individuals besides the service user and that the keys are safe from leakage). In these cases an IAM role should be preferred over an IAM user, if possible.
### Usage ###
```
Usage: main [OPTION] [ARGUMENT]
Used to manage AWS access keys on local system.

	help,				shows program help
	rotate [PROFILE_NAME],		rotates access keys for [PROFILE_NAME]
			all,		rotates all access keys in '$HOME/.aws/credentials' file
	auth [PROFILE_NAME],		auths mfa for [PROFILE_NAME]
```
