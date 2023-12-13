package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	StandardClaims jwt.RegisteredClaims
	Roles          []string
}

func (c MyClaims) Valid() error {
	return nil
}

func main() {

	err := genKey()
	err = genToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func genToken() error {

	filename := "./leo/keys/529484fe-9989-11ee-b9d1-0242ac120002.pem"

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open private key file: %w", err)
	}
	privatePEM, err := io.ReadAll(io.LimitReader(file, 1024*1024))
	if err != nil {
		return fmt.Errorf("reading auth private: %w", err)
	}

	PrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return fmt.Errorf("parsing auth private key: %w", err)
	}

	// ==============================================================

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)

	claims := MyClaims{
		StandardClaims: jwt.RegisteredClaims{
			Issuer:    "service project",
			Subject:   "12345678",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: []string{"ADMIN"},
	}
	method := jwt.GetSigningMethod("RS256")
	token := jwt.NewWithClaims(method, claims)
	token.Header["kid"] = "529484fe-9989-11ee-b9d1-0242ac120002"

	// ==============================================================

	str, err := token.SignedString(PrivateKey)
	if err != nil {
		return err
	}

	fmt.Println("======= TOKEN BEGIN =============")
	fmt.Println(str)
	fmt.Println("======= TOKEN END ===============")
	fmt.Print("\n")

	// ==============================================================

	//Marshal the public key from private key to PRIX.
	asnBytes1, err := x509.MarshalPKIXPublicKey(&PrivateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asnBytes1,
	}

	// Write public key to the public key file.
	if err := pem.Encode(os.Stdout, &publicBlock); err != nil {
		return fmt.Errorf("encoding to publci file: %w", err)
	}

	return nil
}

// Func for generating an x509 private/public keys for auth tokens
func genKey() error {
	//Generate a bnew private key.
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	//Create a file for the private key information in PEM format.
	privateFile, err := os.Create("./leo/keys/529484fe-9989-11ee-b9d1-0242ac120002.pem")
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
