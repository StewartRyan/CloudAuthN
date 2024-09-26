package microsoft

import (
	"fmt"
	"net/url"
	"time"

	"github.com/stewartryan/cloudauthn/internal/utils"

	"github.com/google/uuid"
)

// GenerateSAMLLoginURL generates a SAML login URL for Microsoft Azure Entra ID (formerly Active Directory).
// It takes three parameters:
// - appIdUri: The application ID URI of the Azure AD application.
// - tenantId: The tenant ID of the Azure AD directory.
// - assertionConsumerServiceURL: The URL where the SAML assertion should be sent after successful authentication.
//
// The function returns a string representing the generated SAML login URL and an error if any occurs during the process.
//
// The generated URL follows the format:
// https://login.microsoftonline.com/{tenantId}/saml2?SAMLRequest={encodedSamlRequest}
//
// The SAML request is constructed using the provided parameters and encoded using the RawDeflateBase64Encode function from the utils package.
func GenerateSAMLLoginURL(appIdUri, tenantId, assertionConsumerServiceURL string) (string, error) {
	id := uuid.New().String()

	samlRequest := fmt.Sprintf(`
        <samlp:AuthnRequest xmlns="urn:oasis:names:tc:SAML:2.0:metadata" ID="id%s" Version="2.0" IssueInstant="%s" IsPassive="false" AssertionConsumerServiceURL="%s" xmlns:samlp="urn:oasis:names:tc:SAML:2.0:protocol">
            <Issuer xmlns="urn:oasis:names:tc:SAML:2.0:assertion">%s</Issuer>
            <samlp:NameIDPolicy Format="urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress"></samlp:NameIDPolicy>
        </samlp:AuthnRequest>
        `, id, time.Now().UTC().Format(time.RFC3339), assertionConsumerServiceURL, appIdUri)

	samlBase64, err := utils.RawDeflateBase64Encode(samlRequest)

	if err != nil {
		return "", err
	}

	loginURL := fmt.Sprintf("https://login.microsoftonline.com/%s/saml2?SAMLRequest=%s", tenantId, url.QueryEscape(samlBase64))
	return loginURL, nil
}
