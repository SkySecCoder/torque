# Torque
### Introduction ###
Torque is a simple and secure tool for managing AWS IAM keys and securely interfacing with other AWS resources.
Features include:
- Rotating AWS access keys
- Using AWS access keys with MFA enforced/enabled
- Assuming an IAM role to another account
- Securely SSH'ing into AWS instances without managing ssh keys(work in progress)
### Objective ###
AWS by deafult allows the use of it's access keys without the user needing to use their MFA.

This is a security flaw especially if the keys are leaked to unauthorized parties.

These can be leaked due to accidental code push to GitHub, or simply by copying the unencrypted keys from the `~/.aws/credentials` file to an unauthorised location.

In order to protect these keys, the IAM user must have a policy that forces them to use MFA while using their keys. This can be achieved by enfrocing MFA via an IAM policy and using `Torque` to create temporary access keys that the user can use for their work.

Unfortunately Torque cannot be used when dealing with service users who may need the key but are unable to provide their MFA.(Note: you should make sure that the service user keys are regularly rotated, are not accessible by many individuals besides the service user and that the keys are safe from leakage). In these cases an IAM role should be preferred over an IAM user, if possible.

Torque gives users a convinient and secure way of managing their AWS keys as well as using these keys to interface with AWS resources
### Usage ###
```
Usage: main [OPTION] [ARGUMENT]
Used to manage AWS access keys on local system.

    help,               shows program help
    rotate [PROFILE_NAME],      rotates access keys for [PROFILE_NAME]
            all,        rotates all access keys in '$HOME/.aws/credentials' file
    auth [PROFILE_NAME],        auths mfa for [PROFILE_NAME]
    assume [PROFILE_NAME],      assumes IAM role as specified in [PROFILE_NAME]
```

In order to assume a role via torque you need to specify credentials and configuration in files such as the shared configuration file (`~/.aws/config`) and the shared credentials file (`~/.aws/credentials`) like so :
```
[profile assume_role_profile]
role_arn = arn:aws:iam::<account_number>:role/<role_name>
source_profile = other_profile
mfa_serial = <hardware device serial number or virtual device arn>
```
To configure your configuration profile to assume an IAM role with MFA, you need to specify the MFA device’s serial number for a Hardware MFA device, or ARN for a Virtual MFA device (mfa_serial). This is in addition to specifying the role’s ARN (role_arn).
- assume_role_profile : the name of the profile that you wish to assume roles into
- role_arn : arn of the role that you wish to assume
- source_profile : the profile you would like to use in order to assume your role
- mfa_serial : your MFA serial

This is a one time setup and provides the user a more convinient and secure way of accessing their keys while being sure that the keys themselves are safe(as they would need their MFA to access them)
