package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

type LicenseClaims struct {
	ClientID   string `json:"client_id"`
	ClientName string `json:"client_name"`
	HardwareID string `json:"hardware_id"`
	TrialLimit int    `json:"trial_limit"`
	jwt.RegisteredClaims
}

// InitRSAKeys loads RSA keys or generates them if they don't exist
func InitRSAKeys() error {
	privPath := os.Getenv("RSA_PRIVATE_KEY_PATH")
	if privPath == "" {
		privPath = "private_key.pem"
	}
	pubPath := os.Getenv("RSA_PUBLIC_KEY_PATH")
	if pubPath == "" {
		pubPath = "public_key.pem"
	}

	// 1. Try to load from env variable direct content
	privKeyEnv := os.Getenv("RSA_PRIVATE_KEY")
	if privKeyEnv != "" {
		block, _ := pem.Decode([]byte(privKeyEnv))
		if block == nil {
			return errors.New("failed to parse PEM block from RSA_PRIVATE_KEY env")
		}
		var err error
		PrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			// Try PKCS8
			parsedKey, err8 := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err8 != nil {
				return fmt.Errorf("failed to parse private key: PKCS1=%v, PKCS8=%v", err, err8)
			}
			var ok bool
			PrivateKey, ok = parsedKey.(*rsa.PrivateKey)
			if !ok {
				return errors.New("key is not an RSA private key")
			}
		}
		PublicKey = &PrivateKey.PublicKey
		fmt.Println("RSA private key loaded successfully from environment variable.")
		return nil
	}

	// 2. Try to load from physical files
	if _, err := os.Stat(privPath); os.IsNotExist(err) {
		fmt.Printf("RSA private key not found at %s. Generating a new RSA-2048 keypair...\n", privPath)
		errGen := GenerateRSAKeyPair(privPath, pubPath)
		if errGen != nil {
			return fmt.Errorf("failed to generate RSA keypair: %v", errGen)
		}
	}

	// Load private key file
	privBytes, err := ioutil.ReadFile(privPath)
	if err != nil {
		return fmt.Errorf("failed to read private key file: %v", err)
	}

	block, _ := pem.Decode(privBytes)
	if block == nil {
		return errors.New("failed to parse PEM block from private key file")
	}

	var errParse error
	PrivateKey, errParse = x509.ParsePKCS1PrivateKey(block.Bytes)
	if errParse != nil {
		parsedKey, err8 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err8 != nil {
			return fmt.Errorf("failed to parse private key: PKCS1=%v, PKCS8=%v", errParse, err8)
		}
		var ok bool
		PrivateKey, ok = parsedKey.(*rsa.PrivateKey)
		if !ok {
			return errors.New("key is not an RSA private key")
		}
	}

	// Load public key file
	pubBytes, err := ioutil.ReadFile(pubPath)
	if err == nil {
		pubBlock, _ := pem.Decode(pubBytes)
		if pubBlock != nil {
			pubKeyParsed, errPub := x509.ParsePKIXPublicKey(pubBlock.Bytes)
			if errPub == nil {
				if rsaPub, ok := pubKeyParsed.(*rsa.PublicKey); ok {
					PublicKey = rsaPub
				}
			}
		}
	}

	if PublicKey == nil {
		PublicKey = &PrivateKey.PublicKey
	}

	fmt.Println("RSA-2048 Keypair loaded successfully from disk.")
	return nil
}

// GenerateRSAKeyPair generates an RSA-2048 keypair and writes PEMs to disk
func GenerateRSAKeyPair(privPath, pubPath string) error {
	// Create directory if not exists
	privDir := filepath.Dir(privPath)
	if privDir != "." && privDir != "" {
		_ = os.MkdirAll(privDir, 0755)
	}
	pubDir := filepath.Dir(pubPath)
	if pubDir != "." && pubDir != "" {
		_ = os.MkdirAll(pubDir, 0755)
	}

	// Generate key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// Save private key
	privFile, err := os.Create(privPath)
	if err != nil {
		return err
	}
	defer privFile.Close()

	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	err = pem.Encode(privFile, privBlock)
	if err != nil {
		return err
	}

	// Save public key
	pubFile, err := os.Create(pubPath)
	if err != nil {
		return err
	}
	defer pubFile.Close()

	pubASN1, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	}
	err = pem.Encode(pubFile, pubBlock)
	if err != nil {
		return err
	}

	fmt.Printf("Generated %s and %s\n", privPath, pubPath)
	return nil
}

// SignLicenseToken generates a signed JWT license token
func SignLicenseToken(clientID, clientName, hardwareID string, trialLimit int, expiresAt *time.Time) (string, error) {
	if PrivateKey == nil {
		return "", errors.New("RSA Private Key is not initialized")
	}

	var expTime *jwt.NumericDate
	if expiresAt != nil {
		expTime = jwt.NewNumericDate(*expiresAt)
	}

	claims := LicenseClaims{
		ClientID:   clientID,
		ClientName: clientName,
		HardwareID: hardwareID,
		TrialLimit: trialLimit,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "churma-keygen",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: expTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

// GetPublicKeyPEM returns the public key in PEM format
func GetPublicKeyPEM() (string, error) {
	if PublicKey == nil {
		return "", errors.New("RSA Public Key is not initialized")
	}

	pubASN1, err := x509.MarshalPKIXPublicKey(PublicKey)
	if err != nil {
		return "", err
	}

	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	}

	pemBytes := pem.EncodeToMemory(pubBlock)
	return string(pemBytes), nil
}
