# CloudAuthN

CloudAuthN is a Go-based command line tool designed to retrieve AWS credentials by authenticating through an Identity Provider (IdP) Single Sign-On (SSO). Currently, only Microsoft IdP and AWS are supported, but the tool is designed with extensibility in mind for future cloud providers and IdPs.

## Features

- **AWS Credential Retrieval**: Automatically fetch all AWS credentials for which the user has been granted access, after successful authentication.
- **Microsoft SSO Authentication**: Authenticate using Microsoft as the Identity Provider (IdP) via SSO.
- **Command Line Interface**: Simple and efficient CLI interface built using Golang.

## Requirements

- **Go 1.20+**: Make sure you have Go installed on your machine.
- This tool currently authenticates via Microsoft SSO to Authenticate against AWS via SAML. This means you will need to configure the following environment variables as follows:
  - `CLOUDAUTHN_MICROSOFT_TENANT_ID` — This is the Microsoft Azure Tenant which is used to authenticate your corporation or organisation's users. This ID can be found in the login URL `https://login.microsoftonline.com/{tenantId}`

## Installation

Clone the repository:
```bash
git clone https://github.com/StewartRyan/CloudAuthN.git
```

Navigate into the project directory:
```bash
cd cloudauthn
```

Build the CLI:
```bash
go build -o cloudauthn
```

Move the binary to a directory in your $PATH, e.g.:
```bash
mv cloudauthn /usr/local/bin/
```

## Usage

To retrieve AWS credentials, simply run the following command:

```bash
cloudauthn
```

## Contributing

1. Find an issue or open a new one.
2. Fork the project and create a new branch.
3. Implement changes, test thoroughly, and commit.
4. Open a pull request, address feedback, and merge when approved.


Thank you for your interest in contributing to this project!

## Roadmap

- Built-in quick AWS profile switching
- Support for additional cloud providers (Azure, GCP).
- Additional IdP support (Okta, Google Workspaces).

## License

This project is licensed under the Apache License Version 2.0.