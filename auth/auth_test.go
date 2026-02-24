package auth_test

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/amiwrpremium/go-thalex/apierr"
	"github.com/amiwrpremium/go-thalex/auth"
)

// sha512Sum is a helper that computes SHA-512 and returns the hash slice.
func sha512Sum(data []byte) []byte {
	h := sha512.Sum512(data)
	return h[:]
}

// generateTestKey creates an RSA key pair for testing purposes.
func generateTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate test RSA key: %v", err)
	}
	return key
}

// encodePKCS1PEM encodes a private key in PKCS1 PEM format.
func encodePKCS1PEM(key *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}

// encodePKCS8PEM encodes a private key in PKCS8 PEM format.
func encodePKCS8PEM(t *testing.T, key *rsa.PrivateKey) []byte {
	t.Helper()
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal PKCS8 key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	})
}

// generateECKey creates an ECDSA key for testing non-RSA PKCS8 paths.
func generateECKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func TestNewCredentialsFromPEM_PKCS1(t *testing.T) {
	key := generateTestKey(t)
	pemData := encodePKCS1PEM(key)

	creds, err := auth.NewCredentialsFromPEM("test-key-id", pemData)
	if err != nil {
		t.Fatalf("NewCredentialsFromPEM() unexpected error: %v", err)
	}
	if creds.KeyID != "test-key-id" {
		t.Errorf("KeyID = %q, want %q", creds.KeyID, "test-key-id")
	}
	if creds.PrivateKey == nil {
		t.Fatal("PrivateKey should not be nil")
	}
	if creds.PrivateKey.N.Cmp(key.N) != 0 {
		t.Error("PrivateKey does not match the original key")
	}
}

func TestNewCredentialsFromPEM_PKCS8(t *testing.T) {
	key := generateTestKey(t)
	pemData := encodePKCS8PEM(t, key)

	creds, err := auth.NewCredentialsFromPEM("pkcs8-key", pemData)
	if err != nil {
		t.Fatalf("NewCredentialsFromPEM() unexpected error: %v", err)
	}
	if creds.KeyID != "pkcs8-key" {
		t.Errorf("KeyID = %q, want %q", creds.KeyID, "pkcs8-key")
	}
	if creds.PrivateKey == nil {
		t.Fatal("PrivateKey should not be nil")
	}
	if creds.PrivateKey.N.Cmp(key.N) != 0 {
		t.Error("PrivateKey does not match the original key")
	}
}

func TestNewCredentialsFromPEM_NoPEMBlock(t *testing.T) {
	_, err := auth.NewCredentialsFromPEM("key-id", []byte("not a PEM block"))
	if err == nil {
		t.Fatal("expected error for invalid PEM data, got nil")
	}
	var authErr *apierr.AuthError
	if !errors.As(err, &authErr) {
		t.Errorf("expected AuthError, got %T", err)
	}
}

func TestNewCredentialsFromPEM_WrongBlockType(t *testing.T) {
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte("fake data"),
	})

	_, err := auth.NewCredentialsFromPEM("key-id", pemData)
	if err == nil {
		t.Fatal("expected error for wrong PEM block type, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported PEM block type") {
		t.Errorf("error should mention unsupported block type, got: %v", err)
	}
	var authErr *apierr.AuthError
	if !errors.As(err, &authErr) {
		t.Errorf("expected AuthError, got %T", err)
	}
}

func TestNewCredentialsFromPEM_CorruptPKCS1Data(t *testing.T) {
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: []byte("corrupt key data"),
	})

	_, err := auth.NewCredentialsFromPEM("key-id", pemData)
	if err == nil {
		t.Fatal("expected error for corrupt key data, got nil")
	}
	var authErr *apierr.AuthError
	if !errors.As(err, &authErr) {
		t.Errorf("expected AuthError, got %T", err)
	}
}

func TestNewCredentialsFromPEM_CorruptPKCS8Data(t *testing.T) {
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: []byte("corrupt key data"),
	})

	_, err := auth.NewCredentialsFromPEM("key-id", pemData)
	if err == nil {
		t.Fatal("expected error for corrupt PKCS8 key data, got nil")
	}
	var authErr *apierr.AuthError
	if !errors.As(err, &authErr) {
		t.Errorf("expected AuthError, got %T", err)
	}
	if authErr.Err == nil {
		t.Error("CorruptPKCS8 AuthError should wrap an underlying error")
	}
}

func TestNewCredentialsFromPEM_PKCS8NonRSAKey(t *testing.T) {
	ecKey, err := generateECKey()
	if err != nil {
		t.Fatalf("failed to generate EC key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(ecKey)
	if err != nil {
		t.Fatalf("failed to marshal EC key to PKCS8: %v", err)
	}
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	})

	_, credErr := auth.NewCredentialsFromPEM("key-id", pemData)
	if credErr == nil {
		t.Fatal("expected error for non-RSA PKCS8 key, got nil")
	}
	if !strings.Contains(credErr.Error(), "not an RSA key") {
		t.Errorf("error should mention non-RSA key, got: %v", credErr)
	}
	var authErr *apierr.AuthError
	if !errors.As(credErr, &authErr) {
		t.Errorf("expected AuthError, got %T", credErr)
	}
}

func TestNewCredentials(t *testing.T) {
	key := generateTestKey(t)
	creds := auth.NewCredentials("my-key", key)

	if creds.KeyID != "my-key" {
		t.Errorf("KeyID = %q, want %q", creds.KeyID, "my-key")
	}
	if creds.PrivateKey != key {
		t.Error("PrivateKey should be the exact same pointer")
	}
}

func TestNewCredentials_NilKey(t *testing.T) {
	creds := auth.NewCredentials("my-key", nil)
	if creds.KeyID != "my-key" {
		t.Errorf("KeyID = %q, want %q", creds.KeyID, "my-key")
	}
	if creds.PrivateKey != nil {
		t.Error("PrivateKey should be nil when nil is passed")
	}
}

func TestGenerateToken_ValidJWTStructure(t *testing.T) {
	key := generateTestKey(t)
	creds := auth.NewCredentials("test-kid", key)

	token, err := creds.GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken() unexpected error: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("JWT should have 3 dot-separated parts, got %d", len(parts))
	}

	// Each part should be valid base64url
	for i, part := range parts {
		if part == "" {
			t.Errorf("JWT part %d is empty", i)
		}
		_, err := base64.RawURLEncoding.DecodeString(part)
		if err != nil {
			t.Errorf("JWT part %d is not valid base64url: %v", i, err)
		}
	}
}

func TestGenerateToken_HeaderContents(t *testing.T) {
	key := generateTestKey(t)
	creds := auth.NewCredentials("my-kid-123", key)

	token, err := creds.GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken() unexpected error: %v", err)
	}

	parts := strings.Split(token, ".")
	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("failed to decode header: %v", err)
	}

	var header map[string]string
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		t.Fatalf("failed to unmarshal header JSON: %v", err)
	}

	if header["alg"] != "RS512" {
		t.Errorf("header alg = %q, want %q", header["alg"], "RS512")
	}
	if header["typ"] != "JWT" {
		t.Errorf("header typ = %q, want %q", header["typ"], "JWT")
	}
	if header["kid"] != "my-kid-123" {
		t.Errorf("header kid = %q, want %q", header["kid"], "my-kid-123")
	}
}

func TestGenerateToken_PayloadContainsIAT(t *testing.T) {
	key := generateTestKey(t)
	creds := auth.NewCredentials("kid", key)

	token, err := creds.GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken() unexpected error: %v", err)
	}

	parts := strings.Split(token, ".")
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		t.Fatalf("failed to unmarshal payload JSON: %v", err)
	}

	iat, ok := payload["iat"]
	if !ok {
		t.Fatal("payload should contain 'iat' claim")
	}
	iatFloat, ok := iat.(float64)
	if !ok {
		t.Fatalf("iat should be a number, got %T", iat)
	}
	if iatFloat <= 0 {
		t.Errorf("iat should be positive, got %v", iatFloat)
	}
}

func TestGenerateToken_NilPrivateKey(t *testing.T) {
	creds := auth.NewCredentials("kid", nil)

	_, err := creds.GenerateToken()
	if err == nil {
		t.Fatal("expected error for nil private key, got nil")
	}
	var authErr *apierr.AuthError
	if !errors.As(err, &authErr) {
		t.Errorf("expected AuthError, got %T", err)
	}
	if !strings.Contains(err.Error(), "private key is nil") {
		t.Errorf("error should mention nil private key, got: %v", err)
	}
}

func TestGenerateToken_DifferentTokensPerCall(t *testing.T) {
	key := generateTestKey(t)
	creds := auth.NewCredentials("kid", key)

	token1, err := creds.GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken() call 1 error: %v", err)
	}
	token2, err := creds.GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken() call 2 error: %v", err)
	}

	// Both should have valid structure
	if len(strings.Split(token1, ".")) != 3 {
		t.Error("token1 should have 3 parts")
	}
	if len(strings.Split(token2, ".")) != 3 {
		t.Error("token2 should have 3 parts")
	}
}

// TestGenerateToken_SigningFailure triggers the rsa.SignPKCS1v15 error path
// by using an RSA key with a tiny modulus that is too small for SHA-512 signing.
func TestGenerateToken_SigningFailure(t *testing.T) {
	// Construct a deliberately broken RSA key with a tiny modulus.
	// This triggers the signing error because PKCS1v15 requires the
	// modulus to be large enough to hold the hash prefix + digest.
	smallKey := &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: new(big.Int).SetInt64(3233), // tiny modulus
			E: 17,
		},
		D: new(big.Int).SetInt64(2753),
	}

	creds := auth.NewCredentials("small-key", smallKey)
	_, tokenErr := creds.GenerateToken()
	if tokenErr == nil {
		t.Fatal("expected error when signing with a key too small for RS512, got nil")
	}

	var authErr *apierr.AuthError
	if !errors.As(tokenErr, &authErr) {
		t.Errorf("expected AuthError, got %T: %v", tokenErr, tokenErr)
	}
	if !strings.Contains(tokenErr.Error(), "failed to sign JWT") {
		t.Errorf("error should mention signing failure, got: %v", tokenErr)
	}
	if authErr.Err == nil {
		t.Error("signing failure AuthError should wrap an underlying error")
	}
}

// TestGenerateToken_VerifySignature verifies the JWT signature is valid
// by checking it against the public key.
func TestGenerateToken_VerifySignature(t *testing.T) {
	key := generateTestKey(t)
	creds := auth.NewCredentials("verify-kid", key)

	token, err := creds.GenerateToken()
	if err != nil {
		t.Fatalf("GenerateToken() unexpected error: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("JWT should have 3 parts, got %d", len(parts))
	}

	signingInput := parts[0] + "." + parts[1]
	sigBytes, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		t.Fatalf("failed to decode signature: %v", err)
	}

	// Verify signature using the public key
	hash := sha512Sum([]byte(signingInput))
	verifyErr := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA512, hash, sigBytes)
	if verifyErr != nil {
		t.Errorf("signature verification failed: %v", verifyErr)
	}
}

// TestNewCredentialsFromPEM_EmptyKeyID ensures empty key ID is accepted.
func TestNewCredentialsFromPEM_EmptyKeyID(t *testing.T) {
	key := generateTestKey(t)
	pemData := encodePKCS1PEM(key)

	creds, err := auth.NewCredentialsFromPEM("", pemData)
	if err != nil {
		t.Fatalf("NewCredentialsFromPEM() unexpected error: %v", err)
	}
	if creds.KeyID != "" {
		t.Errorf("KeyID = %q, want empty string", creds.KeyID)
	}
}

// TestNewCredentialsFromPEM_EmptyPEMData ensures empty data returns an error.
func TestNewCredentialsFromPEM_EmptyPEMData(t *testing.T) {
	_, err := auth.NewCredentialsFromPEM("key-id", []byte{})
	if err == nil {
		t.Fatal("expected error for empty PEM data, got nil")
	}
	var authErr *apierr.AuthError
	if !errors.As(err, &authErr) {
		t.Errorf("expected AuthError, got %T", err)
	}
	if !strings.Contains(err.Error(), "failed to decode PEM block") {
		t.Errorf("error should mention PEM decode failure, got: %v", err)
	}
}

// TestNewCredentialsFromPEM_NilPEMData ensures nil data returns an error.
func TestNewCredentialsFromPEM_NilPEMData(t *testing.T) {
	_, err := auth.NewCredentialsFromPEM("key-id", nil)
	if err == nil {
		t.Fatal("expected error for nil PEM data, got nil")
	}
	var authErr *apierr.AuthError
	if !errors.As(err, &authErr) {
		t.Errorf("expected AuthError, got %T", err)
	}
}

// TestNewCredentialsFromPEM_CorruptPKCS1_WrapsError ensures the corrupt PKCS1
// path wraps the underlying parse error.
func TestNewCredentialsFromPEM_CorruptPKCS1_WrapsError(t *testing.T) {
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: []byte("corrupt key data"),
	})

	_, err := auth.NewCredentialsFromPEM("key-id", pemData)
	if err == nil {
		t.Fatal("expected error for corrupt key data, got nil")
	}
	var authErr *apierr.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthError, got %T", err)
	}
	if !strings.Contains(err.Error(), "failed to parse RSA private key") {
		t.Errorf("error should mention parse failure, got: %v", err)
	}
	if authErr.Err == nil {
		t.Error("corrupt PKCS1 AuthError should wrap an underlying error")
	}
}
