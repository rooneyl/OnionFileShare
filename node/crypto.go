package node

import ()

// type of PublicKey and PrivateKey should be change
// i just set that as string to make a stub

func encryptData(data []byte, publicKey string) ([]byte, error) {
	var encryptedData []byte
	return encryptedData, nil
}

func decryptData(data []byte, PrivateKey string) ([]byte, error) {
	var decryptedData []byte
	return decryptedData, nil
}

func generateKeys() (string, string) {
	return "publicKey", "privateKey"
}
