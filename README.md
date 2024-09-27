# CloudAuthN

CloudAuthN is a Go-based command line tool designed to retrieve AWS credentials by authenticating through an Identity Provider (IdP) Single Sign-On (SSO). Currently, only Microsoft IdP and AWS are supported, but the tool is designed with extensibility in mind for future cloud providers and IdPs.

## Features

- **AWS Credential Retrieval**: Automatically fetch AWS credentials after successful authentication.
- **Microsoft SSO Authentication**: Authenticate using Microsoft as the Identity Provider (IdP) via SSO.
- **Command Line Interface**: Simple and efficient CLI interface built using Golang.

## Requirements

- **Go 1.20+**: Make sure you have Go installed on your machine.
- **AWS Account**: You must have access to AWS services with proper permissions.
- **Microsoft SSO Account**: Authentication is done through Microsoft SSO.

## Installation

Clone the repository:
```bash
git clone https://github.com/StewartRyan/CloudAuthN.git
```

Navigate into the project directory:
```bash
cd cloud-credentials-cli
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

To retrieve AWS credentials, run the following command:

```bash
cloudauthn --provider aws --idp microsoft --role <your-aws-role> --profile <aws-profile-name>
```

Example:

```bash
cloudauthn --provider aws --idp microsoft --role admin --profile default
```

## Flags

- `--provider`: Specify the cloud provider. Currently, only `aws` is supported.
- `--idp`: Specify the Identity Provider. Currently, only `microsoft` is supported.

## Roadmap

- Support for additional cloud providers (Azure, GCP).
- Additional IdP support (Okta, Google).
- Token caching for faster reauthentication.

## License

This project is licensed under the Apache License Version 2.0.