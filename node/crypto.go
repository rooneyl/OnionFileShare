package node

import ()

// type of PublicKey and PrivateKey should be change
// i just set that as string to make a stub

func encryptData(data []byte, publicKey string) ([]byte, error) {
	var encryptedData []byte
	return encryptedData, nil
}

func decryptData(data []byte, privateKey string) ([]byte, error) {
	var decryptedData []byte
	return decryptedData, nil
}

func encryptStruct(stru interface{}, publicKey string) ([]byte, error) {
	var encryptedStruct []byte
	return encryptedStruct, nil
}

func decryptStruct(data []byte, privateKey string, stru *interface{}) error {
	return nil
}

func generateKeys() (string, string) {
	return "publicKey", "privateKey"
}
