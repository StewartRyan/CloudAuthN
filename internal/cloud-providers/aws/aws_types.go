package aws

import "time"

type AwsRole struct {
	RoleArn      string
	PrincipalArn string
}

type AWSResponse struct {
	Roles        []AwsRole
	SAMLResponse string
}

type AwsCredentialsEntry struct {
	AccessKeyID     string    `yaml:"aws_access_key_id"`
	SecretAccessKey string    `yaml:"aws_secret_access_key"`
	SessionToken    string    `yaml:"aws_session_token"`
	Expiration      time.Time `yaml:"expiration"`
}
