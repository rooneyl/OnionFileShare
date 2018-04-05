package node

import (
	"testing"
)

type TestStruct struct {
	TestField string
}

func TestDecryptData(t *testing.T) {
	struc1 := TestStruct{"this is my test field"}
	var struc2 TestStruct
	aesKey := generateAESKey()
	encData, err := encryptData(aesKey, &struc1)

	if err != nil {
		t.Errorf("Error encrypting the data, Got: %v, Expected: %v.", err, nil)
	}

	err = decryptData(aesKey,encData,&struc2)
	if err != nil {
		t.Errorf("Error decrypting the data, Got: %v, Expected: %v.", err, nil)

	}

	if struc1.TestField != struc2.TestField {
		t.Errorf("Struc1 and Struc2 are not equivalent, Got: %s, Expected: %s.", struc2.TestField, struc1.TestField)
	}
	// Output: MOOOO!
}

func TestDecryptAESKey(t *testing.T) {
	aesKey := generateAESKey();
	pubBytes, privBytes := generateRSAKey()

	encAESKey, err := encryptAESKey(aesKey, pubBytes)
	if err != nil {
		t.Errorf("Sum was incorrect, Got: %v, Expected: %v.", err, nil)

	}
	aesKey2, err := decryptAESKey(encAESKey, privBytes);

	if string(aesKey) != string(aesKey2) {
		t.Errorf("EASKeys not equivalent, Got: %s, Expected: %s.", aesKey2, aesKey)
	}
	// Output: MOOOO!
}

