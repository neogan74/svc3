package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	err := genKey()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Func for generating an x509 private/public keys for auth tokens
func genKey() error {
	//Generate a bnew private key.
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	//Create a file for the private key information in PEM format.
	privateFile, err := os.Create("private.pem")
	if err != nil {
		return fmt.Errorf("Error with creating private file: %w", err)
	}
	defer privateFile.Close()

	//Construct a PEM block for the private key
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}

	//Write the private key to the private file.
	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return fmt.Errorf("encoding to private file: %w", err)
	}

	// =================================================================

	//Marshal the public key from private key to PRIX.
	asnBytes1, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	publicFile, err := os.Create("public.pem")
	if err != nil {
		return fmt.Errorf("Creating public file: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asnBytes1,
	}

	// Write public key to the public key file.
	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding to publci file: %w", err)
	}

	fmt.Println("Private and publc key were geenrated.")

	return nil
}
