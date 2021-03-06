package node

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"crypto/aes"
	"io"
	"crypto/cipher"
	"errors"
)

func generateRSAKey() ([]byte, []byte) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	checkError(err)

	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
	pubBytes := x509.MarshalPKCS1PublicKey(&privKey.PublicKey)
	return pubBytes, privBytes
}

func generateAESKey() []byte {
	rng := rand.Reader
	key := make([]byte, 32)
	if _, err := io.ReadFull(rng, key); err != nil {
		Log.Fatal(err)
	}
	return key
}

func encryptAESKey(aesKey []byte, pubBytes []byte) ([]byte, error) {
	label := []byte("orders")

	pubKey, err := x509.ParsePKCS1PublicKey(pubBytes)
	checkError(err)
	//rsa.EncryptPKCS1v15(rng,pubKey,cipherText)
	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, aesKey, label)
	checkError(err)
	return encryptedAESKey, err
}

func decryptAESKey(encAESKey []byte, privBytes []byte) ([]byte, error) {

	privKey, err := x509.ParsePKCS1PrivateKey(privBytes)
	checkError(err)

	//decrypt aesKey using private key and sha256 algorithm
	label := []byte("orders")
	aesKey, err :=rsa.DecryptOAEP(sha256.New(),rand.Reader, privKey, encAESKey, label)
	return aesKey, err
}

func encryptData(aesKey []byte, struc interface{}) ([]byte, error) {
	data, err := json.Marshal(struc)
	checkError(err)

	block, err := aes.NewCipher(aesKey)
	checkError(err)

	encryptedData := make([]byte, aes.BlockSize+len(data))
	iv := encryptedData[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		Log.Fatal(err)
	}

	stream := cipher.NewCFBEncrypter(block,iv)
	stream.XORKeyStream(encryptedData[aes.BlockSize:], data)

	return encryptedData, err
}

func decryptData(aesKey []byte, data []byte, struc interface{}) error {
	//make a new block using EASKey
	block, err := aes.NewCipher(aesKey)
	checkError(err)

	//make sure blocksize and data are the same length
	if len(data) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		Log.Fatal(err)
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	//create a new decryptor with the block and IV
	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	//decrypt data
	stream.XORKeyStream(data, data)

	json.Unmarshal(data, struc)
	checkError(err)

	return nil
}

//func EncryptStruct(struc interface{}, pubBytes []byte) ([]byte, []byte, error) {
//
//	data, err := json.Marshal(struc)
//	checkError(err)
//
//	//generate a symmetric key for AES encryption (session key)
//	rng := rand.Reader
//	key := make([]byte, 32)
//	if _, err := io.ReadFull(rng, key); err != nil {
//		Log.Fatal(err)
//	}
//
//	//create a new block using the AES key
//	block, err := aes.NewCipher(key)
//	checkError(err)
//
//	pubKey, err := x509.ParsePKCS1PublicKey(pubBytes)
//	checkError(err)
//
//	encryptedData := make([]byte, aes.BlockSize+len(data))
//	iv := encryptedData[:aes.BlockSize]
//
//	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
//	Log.Fatal(err)
//	}
//
//	stream := cipher.NewCFBEncrypter(block,iv)
//	stream.XORKeyStream(encryptedData[aes.BlockSize:], data)
//
//	//rsaPubKey := pubKey.(*rsa.PublicKey)
//	label := []byte("orders")
//
//	//rsa.EncryptPKCS1v15(rng,pubKey,cipherText)
//	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, key, label)
//	checkError(err)
//
//	return encryptedAESKey, encryptedData , nil
//}
//
//func DecryptStruct(data []byte, encryptedEosKey []byte, priByte []byte, struc interface{}) error {
//
//	//get private key from private key bytes
//	privKey, err := x509.ParsePKCS1PrivateKey(priByte)
//	if err != nil {
//		Log.Fatal(err)
//	}
//
//	//decrypt aesKey using private key and sha256 algorithm
//	label := []byte("orders")
//	aesKey, err :=rsa.DecryptOAEP(sha256.New(),rand.Reader, privKey, encryptedEosKey, label)
//
//	//make a new block using EASKey
//	block, err := aes.NewCipher(aesKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	//make sure blocksize and data are the same length
//	if len(data) < aes.BlockSize {
//		err = errors.New("Ciphertext block size is too short!")
//		Log.Fatal(err)
//	}
//
//	iv := data[:aes.BlockSize]
//	data = data[aes.BlockSize:]
//
//	//create a new decryptor with the block and IV
//	stream := cipher.NewCFBDecrypter(block, iv)
//
//	// XORKeyStream can work in-place if the two arguments are the same.
//	//decrypt data
//	stream.XORKeyStream(data, data)
//
//	//unmarshal json data into stru
//	json.Unmarshal(data, struc)
//	if err != nil {
//		Log.Fatal(err)
//	}
//
//	return nil
//}

//Generate public and private rsa keys
//func GenerateKeys() ([]byte, []byte) {
//
//	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
//	checkError(err)
//
//	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
//	pubBytes := x509.MarshalPKCS1PublicKey(&privKey.PublicKey)
//
//	return pubBytes, privBytes
//}

func checkError(err error) {
	if err != nil {
		Log.Fatal(err)
	}
}
