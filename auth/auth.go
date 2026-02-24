// Package auth provides API authentication for the Thalex SDK.
//
// It supports RSA key-based authentication using RS512-signed JWT tokens.
// Keys can be loaded from PEM-encoded data in either PKCS1 or PKCS8 format.
package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/amiwrpremium/go-thalex/apierr"
)

// Credentials holds the API key ID and private key for authentication.
type Credentials struct {
	// KeyID is the API key identifier (kid).
	KeyID string
	// PrivateKey is the RSA private key used for signing JWT tokens.
	PrivateKey *rsa.PrivateKey
}

// NewCredentialsFromPEM creates Credentials from a PEM-encoded RSA private key.
// It supports both PKCS1 and PKCS8 PEM formats.
func NewCredentialsFromPEM(keyID string, pemData []byte) (*Credentials, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, &apierr.AuthError{Message: "failed to decode PEM block"}
	}

	var privKey *rsa.PrivateKey
	var err error

	switch block.Type {
	case "RSA PRIVATE KEY":
		privKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
		if parseErr != nil {
			return nil, &apierr.AuthError{Message: "failed to parse PKCS8 private key", Err: parseErr}
		}
		var ok bool
		privKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, &apierr.AuthError{Message: "PKCS8 key is not an RSA key"}
		}
	default:
		return nil, &apierr.AuthError{Message: fmt.Sprintf("unsupported PEM block type: %s", block.Type)}
	}

	if err != nil {
		return nil, &apierr.AuthError{Message: "failed to parse RSA private key", Err: err}
	}

	return &Credentials{KeyID: keyID, PrivateKey: privKey}, nil
}

// NewCredentials creates Credentials from a pre-parsed RSA private key.
func NewCredentials(keyID string, privateKey *rsa.PrivateKey) *Credentials {
	return &Credentials{KeyID: keyID, PrivateKey: privateKey}
}

// GenerateToken creates a signed JWT token using RS512 for API authentication.
// The token contains the key ID (kid) in the header and issued-at (iat) in the payload.
func (c *Credentials) GenerateToken() (string, error) {
	if c.PrivateKey == nil {
		return "", &apierr.AuthError{Message: "private key is nil"}
	}

	// Header
	header := map[string]string{
		"alg": "RS512",
		"typ": "JWT",
		"kid": c.KeyID,
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", &apierr.AuthError{Message: "failed to marshal JWT header", Err: err}
	}

	// Payload
	payload := map[string]any{
		"iat": time.Now().Unix(),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", &apierr.AuthError{Message: "failed to marshal JWT payload", Err: err}
	}

	// Encode header and payload
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signingInput := headerB64 + "." + payloadB64

	// Sign with RS512
	hash := sha512.Sum512([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, c.PrivateKey, crypto.SHA512, hash[:])
	if err != nil {
		return "", &apierr.AuthError{Message: "failed to sign JWT", Err: err}
	}

	sigB64 := base64.RawURLEncoding.EncodeToString(sig)
	return signingInput + "." + sigB64, nil
}
