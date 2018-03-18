package node

import (
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"log"
	"encoding/hex"
	"crypto/sha256"
	"encoding/json"
)

// type of PublicKey and PrivateKey should be change
// i just set that as string to make a stub

func encryptData(data []byte, publicKey string) ([]byte, error) {

	pubKeyBytes, err := hex.DecodeString(publicKey)
	checkError(err)

	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	checkError(err)

	rsaPubKey := pubKey.(*rsa.PublicKey)
	label := []byte("orders")

	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, data, label)
	checkError(err)

	return encryptedData, nil
}

func decryptData(data []byte, privateKey string) ([]byte, error) {

	privKeyBytes, err := hex.DecodeString(privateKey)
	checkError(err)

	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	checkError(err)

	label := []byte("orders")

	decryptedData, err := rsa.DecryptOAEP(sha256.New(),rand.Reader,privKey,data, label)
	checkError(err)

	return decryptedData, nil
}

func encryptStruct(struc interface{}, publicKey string) ([]byte, error) {

	data, err := json.Marshal(struc)
	checkError(err)

	pubKeyBytes, err := hex.DecodeString(publicKey)
	checkError(err)

	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	checkError(err)

	rsaPubKey := pubKey.(*rsa.PublicKey)
	label := []byte("orders")

	encryptedStruct, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, data, label)
	checkError(err)

	return encryptedStruct, nil
}

func decryptStruct(data []byte, privateKey string, stru interface{}) error {
	privKeyBytes, err := hex.DecodeString(privateKey)
	checkError(err)

	privKey, err := x509.ParsePKCS1PrivateKey(privKeyBytes)
	checkError(err)

	label := []byte("orders")

	decryptedData, err := rsa.DecryptOAEP(sha256.New(),rand.Reader,privKey,data, label)

	json.Unmarshal(decryptedData, stru)
	checkError(err)

	return nil
}

func generateKeys() (string, string) {

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	checkError(err)

	privKeyBytes := x509.MarshalPKCS1PrivateKey(priv)

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(priv.PublicKey)
	checkError(err)

	pubKeyString := hex.EncodeToString(pubKeyBytes)
	privKeyString := hex.EncodeToString(privKeyBytes)

	return pubKeyString, privKeyString
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
