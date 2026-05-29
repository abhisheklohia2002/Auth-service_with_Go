package common

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strings"

	"example.com/m/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	config     config.Config
	privateKey *rsa.PrivateKey
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func LoadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	switch block.Type {
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		rsaPub, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("public key is not RSA")
		}

		return rsaPub, nil

	case "RSA PUBLIC KEY":
		return x509.ParsePKCS1PublicKey(block.Bytes)

	default:
		return nil, errors.New("unsupported public key type: " + block.Type)
	}
}

func normalizePEM(value string) string {
	return strings.ReplaceAll(value, `\n`, "\n")
}

func LoadRSAPrivateKeyFromEnv(envKey string) (*rsa.PrivateKey, error) {
	raw := os.Getenv(envKey)
	if raw == "" {
		return nil, errors.New(envKey + " is not set")
	}

	raw = normalizePEM(raw)

	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, errors.New("failed to decode private key PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not RSA")
		}
		return rsaKey, nil
	}

	rsaKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return rsaKey, nil
	}

	return nil, errors.New("failed to parse RSA private key")
}

func LoadRSAPublicKeyFromEnv(envKey string) (*rsa.PublicKey, error) {
	raw := os.Getenv(envKey)
	if raw == "" {
		return nil, errors.New(envKey + " is not set")
	}

	raw = normalizePEM(raw)

	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		return nil, errors.New("failed to decode public key PEM")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err == nil {
		rsaKey, ok := pub.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("public key is not RSA")
		}
		return rsaKey, nil
	}

	rsaKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return rsaKey, nil
	}

	return nil, errors.New("failed to parse RSA public key")
}
