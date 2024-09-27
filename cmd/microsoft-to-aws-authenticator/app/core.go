package microsoft_to_aws

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stewartryan/cloudauthn/internal/cloud-providers/aws"
	"github.com/stewartryan/cloudauthn/internal/sso/microsoft"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// MicrosoftToAws automates the process of logging into Microsoft and assuming AWS roles using SAML.
// It opens a Chrome browser with headed mode, navigates to the Microsoft login page, intercepts the SAML response,
// and then updates the AWS credentials file with the assumed roles.
func MicrosoftToAws() {
	// 1. Create a channel to receive the roles
	awsResponseChan := make(chan aws.AWSResponse)

	// 1.1 Get value of env variable CLOUDAUTHN_MICROSOFT_TENANT_ID
	var tenantID string
	if tenantID = os.Getenv("CLOUDAUTHN_MICROSOFT_TENANT_ID"); tenantID == "" {
		log.Fatal("Missing CLOUDAUTHN_MICROSOFT_TENANT_ID environment variable")
	}

	// 2. Create the login URL
	loginUrl, err := microsoft.GenerateSAMLLoginURL("https://signin.aws.amazon.com/saml", tenantID, "https://signin.aws.amazon.com/saml")
	//
	if err != nil {
		fmt.Println("Error creating login URL:", err)
	} else {
		// 3. Open Chrome on Microsoft login page and listen for SAML response
		openChromeOnMicrosoftLogin(loginUrl, awsResponseChan)
	}

	// 4. Wait for the roles to be received
	awsResponse := <-awsResponseChan

	// 5. Close the channel
	close(awsResponseChan)

	// 6. Assume the roles and update the AWS credentials file
	awsCredentials := make(map[string]aws.AwsCredentialsEntry)
	for _, role := range awsResponse.Roles {
		out, err := aws.AssumeRoleWithSAML(role.RoleArn, role.PrincipalArn, awsResponse.SAMLResponse)

		if err != nil {
			fmt.Println("Error assuming role:", err)
		} else {
			// extract account number from role arn
			accountId := strings.Split(role.RoleArn, ":")[4]
			awsCredentials[accountId] = aws.AwsCredentialsEntry{
				AccessKeyID:     *out.Credentials.AccessKeyId,
				SecretAccessKey: *out.Credentials.SecretAccessKey,
				SessionToken:    *out.Credentials.SessionToken,
				Expiration:      *out.Credentials.Expiration,
			}
		}
	}

	// 7. Update the AWS credentials file
	err = aws.UpdateAWSCredentialsFile(awsCredentials)
	if err != nil {
		fmt.Println("Error updating AWS credentials file:", err)
	} else {
		fmt.Println("AWS credentials file updated successfully")
	}
}

// openChromeOnMicrosoftLogin opens a Chrome browser with headed mode and navigates to the provided login URL.
// It then listens for a SAML response from the Microsoft login page and sends the response to the provided channel.
// The function also handles context cancellation and ensures the channel is closed upon completion.
//
// Parameters:
// - loginUrl: The URL of the Microsoft login page to navigate to.
// - awsResponseChan: A channel to send the AWS response data upon receiving a SAML response.
func openChromeOnMicrosoftLogin(loginUrl string, awsResponseChan chan aws.AWSResponse) {
	// Create a new Chrome context with headed mode
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("window-size", "430,680"),
		chromedp.Flag("no-default-browser-check", false),
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
	)

	// Define contexts
	ctx_, cancel_ := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel_()
	ctx, cancel := chromedp.NewContext(ctx_)
	defer cancel()

	// Intercept SAML Response
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		aws.InterceptChromeAuthRequest(ev, awsResponseChan, cancel)
	})

	// Navigate to the login URL (replace with your actual URL)
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(loginUrl),
	)

	if err != nil {
		log.Fatalf("Failed to navigate to login page: %v", err)
	}

	// wait for context to be done, then close the channel
	<-ctx.Done()
}
