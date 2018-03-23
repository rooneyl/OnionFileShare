package node

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
)

func encryptStruct(struc interface{}, pubByte []byte) ([]byte, error) {
	data, err := json.Marshal(struc)
	if err != nil {
		Log.Fatal(err)
	}

	pubKey, err := x509.ParsePKCS1PublicKey(pubByte)
	if err != nil {
		Log.Fatal(err)
	}

	label := []byte("orders")

	encryptedStruct, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, data, label)
	if err != nil {
		Log.Fatal(err)
	}

	return encryptedStruct, nil
}

func decryptStruct(data []byte, priByte []byte, stru interface{}) error {
	privKey, err := x509.ParsePKCS1PrivateKey(priByte)
	if err != nil {
		Log.Fatal(err)
	}

	label := []byte("orders")

	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, data, label)

	json.Unmarshal(decryptedData, stru)
	if err != nil {
		Log.Fatal(err)
	}

	return nil
}

func generateKeys() ([]byte, []byte) {

	pri, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		Log.Fatal(err)
	}

	priByte := x509.MarshalPKCS1PrivateKey(pri)
	pubByte := x509.MarshalPKCS1PublicKey(&pri.PublicKey)

	return pubByte, priByte
}

func checkError(err error) {
	if err != nil {
		Log.Fatal(err)
	}
}
