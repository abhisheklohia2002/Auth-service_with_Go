package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
)

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func base64urlBigInt(n *big.Int) string {
	return base64.RawURLEncoding.EncodeToString(n.Bytes())
}

func base64urlInt(i int) string {
	b := big.NewInt(int64(i)).Bytes()
	return base64.RawURLEncoding.EncodeToString(b)
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

func PublicKeyToJWK(publicKey *rsa.PublicKey, kid string) JWK {
	return JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: kid,
		Alg: "RS256",
		N:   base64urlBigInt(publicKey.N),
		E:   base64urlInt(publicKey.E),
	}
}

func main() {
	publicKey, err := LoadRSAPublicKey("../openSSL/public.pem")
	if err != nil {
		log.Fatal(err)
	}

	jwk := PublicKeyToJWK(publicKey, "access-key-1")

	jwks := JWKS{
		Keys: []JWK{jwk},
	}

	jsonBytes, err := json.MarshalIndent(jwks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBytes))
}
