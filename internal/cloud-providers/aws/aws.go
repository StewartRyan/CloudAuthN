package aws

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/chromedp/cdproto/network"
)

// UpdateAWSCredentialsFile updates the AWS credentials file with the given credentials.
// The function reads the existing credentials file, updates or inserts the given credentials,
// and writes the updated content back to the file.
//
// Parameters:
//   - credentials: A map containing the profile names as keys and AwsCredentialsEntry structs as values.
//     Each AwsCredentialsEntry struct represents the AWS credentials for a specific profile.
//
// Returns:
// - An error if any error occurs during the process, or nil if the operation is successful.
func UpdateAWSCredentialsFile(credentials map[string]AwsCredentialsEntry) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	credsFilePath := filepath.Join(homeDir, ".aws", "credentials")
	fileContent, err := os.ReadFile(credsFilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := strings.Split(string(fileContent), "\n")
	profileEntries := make(map[string][]string)
	currentProfile := ""

	// Read existing profiles and their entries
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentProfile = strings.Trim(line, "[]")
			profileEntries[currentProfile] = []string{}
		} else if currentProfile != "" {
			profileEntries[currentProfile] = append(profileEntries[currentProfile], line)
		}
	}

	// Update or insert the given credentials
	for profile, entry := range credentials {
		profileEntries[profile] = []string{
			fmt.Sprintf("aws_access_key_id=%s", entry.AccessKeyID),
			fmt.Sprintf("aws_secret_access_key=%s", entry.SecretAccessKey),
			fmt.Sprintf("expiration=%s", entry.Expiration.Format(time.RFC3339)),
		}
		if entry.SessionToken != "" {
			profileEntries[profile] = append(profileEntries[profile], fmt.Sprintf("aws_session_token=%s", entry.SessionToken))
		}
	}

	// Generate updated content
	var updatedContent strings.Builder
	for profile, entries := range profileEntries {
		updatedContent.WriteString(fmt.Sprintf("[%s]\n", profile))
		for _, entry := range entries {
			updatedContent.WriteString(fmt.Sprintf("%s\n", entry))
		}
	}

	// Write the updated content back to the file
	if err := os.WriteFile(credsFilePath, []byte(updatedContent.String()), 0600); err != nil {
		return err
	}

	return nil
}

// AssumeRoleWithSAML assumes a role in AWS using a SAML assertion.
// It creates a new session with the default AWS configuration, creates a new STS client,
// and then calls the AssumeRoleWithSAML API with the provided role ARN, principal ARN, and SAML assertion.
// The function returns the AssumeRoleWithSAMLOutput result and any error encountered during the process.
//
// Parameters:
// - roleArn: The ARN of the IAM role to assume.
// - principalArn: The ARN of the IAM principal (user or federated user) that is making the AssumeRoleWithSAML call.
// - samlAssertion: The base64-encoded SAML assertion obtained from the identity provider.
//
// Returns:
// - *sts.AssumeRoleWithSAMLOutput: The result of the AssumeRoleWithSAML API call.
// - error: An error if any error occurs during the process, or nil if the operation is successful.
func AssumeRoleWithSAML(roleArn, principalArn, samlAssertion string) (*sts.AssumeRoleWithSAMLOutput, error) {
	// Create a new session with the default AWS configuration
	sess := session.Must(session.NewSession())

	// Create a new STS client
	svc := sts.New(sess)

	// Assume the role with SAML
	input := &sts.AssumeRoleWithSAMLInput{
		RoleArn:         aws.String(roleArn),
		PrincipalArn:    aws.String(principalArn),
		SAMLAssertion:   aws.String(samlAssertion),
		DurationSeconds: aws.Int64(3600 * 12), // 12 hours
	}

	// Call the AssumeRoleWithSAML API
	result, err := svc.AssumeRoleWithSAML(input)
	if err != nil {
		return nil, fmt.Errorf("failed to assume role with SAML: %w", err)
	}

	return result, nil
}

// ParseRolesFromSamlResponse extracts AWS roles and the SAML response from a given assertion.
// It decodes the assertion, extracts the SAML response, and parses the roles from the SAML response.
//
// Parameters:
// - assertion: A base64-encoded SAML assertion obtained from the identity provider.
//
// Returns:
// - roles: A slice of AwsRole structs representing the extracted AWS roles.
// - samlResponse: A pointer to a string containing the base64-encoded SAML response.
// - error: An error if any error occurs during the process, or nil if the operation is successful.
func ParseRolesFromSamlResponse(assertion string) ([]AwsRole, *string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(assertion)
	if err != nil {
		return nil, nil, err
	}

	samlResponse := string(decodedBytes)

	// url decode the saml response
	samlResponseUriDecoded, err := url.QueryUnescape(samlResponse)
	if err != nil {
		return nil, nil, err
	}

	// take value after SAMLResponse=
	samlBase64Encoded := strings.Split(samlResponseUriDecoded, "SAMLResponse=")[1]

	// base64 decode the saml response
	samlBase64Decoded, err := base64.StdEncoding.DecodeString(samlBase64Encoded)
	if err != nil {
		return nil, nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(samlBase64Decoded))
	if err != nil {
		return nil, nil, err
	}

	var roles []AwsRole
	doc.Find("Attribute[Name='https://aws.amazon.com/SAML/Attributes/Role'] > AttributeValue").Each(func(i int, s *goquery.Selection) {
		roleAndPrincipal := s.Text()
		parts := strings.Split(roleAndPrincipal, ",")

		var roleIdx, principalIdx int
		if strings.Contains(parts[0], ":role/") {
			roleIdx, principalIdx = 0, 1
		} else {
			roleIdx, principalIdx = 1, 0
		}
		roleArn := strings.TrimSpace(parts[roleIdx])
		principalArn := strings.TrimSpace(parts[principalIdx])
		roles = append(roles, AwsRole{RoleArn: roleArn, PrincipalArn: principalArn})
	})

	return roles, &samlBase64Encoded, nil
}

// InterceptChromeAuthRequest intercepts a network request made by Chrome to the AWS sign-in page
// and processes the SAML response. If a valid SAML response is found, it decodes and inflates
// the response, extracts the AWS roles, and sends the response through a channel.
//
// Parameters:
// - ev: An interface representing the network event. It should be of type *network.EventRequestWillBeSent.
// - responseChan: A channel of type AWSResponse where the processed response will be sent.
// - cancel: A context.CancelFunc to cancel the network request if necessary.
func InterceptChromeAuthRequest(ev interface{}, responseChan chan AWSResponse, cancel context.CancelFunc) {
	switch ev := ev.(type) {
	case *network.EventRequestWillBeSent:
		request := ev.Request

		if request.Method == "POST" && strings.Contains(request.URL, "https://signin.aws.amazon.com/saml") && request.PostDataEntries != nil {
			entry := request.PostDataEntries[0]
			if entry.Bytes != "" {
				// decode and inflate the base64 encoded SAML response but ensure to take only the value after the '='
				rolesFound, samlResponse, err := ParseRolesFromSamlResponse(entry.Bytes)
				if err != nil {
					log.Fatalf("Failed to parse SAML response: %v", err)
				}

				cancel()

				awsResponse := AWSResponse{
					Roles:        rolesFound,
					SAMLResponse: *samlResponse,
				}
				responseChan <- awsResponse
			}

		}
	}
}
