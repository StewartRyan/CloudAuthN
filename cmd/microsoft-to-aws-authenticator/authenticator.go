package main

import (
	core "github.com/stewartryan/cloudauthn/cmd/microsoft-to-aws-authenticator/app"
)

func main() {
	// TODO: Implement args parsing and error handling
	// Currently, only Microsoft to AWS authentication is supported
	core.MicrosoftToAws()
}
