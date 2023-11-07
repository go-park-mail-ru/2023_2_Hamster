package hasher

import (
	"fmt"
	"testing"
)

func TestPasswordHashingAndVerification(t *testing.T) {
	// Test GeneratePasswordHash and VerfiyPassword functions

	// Test case 1: Valid password hashing and verification
	password := "securePassword123"
	hash, err := GeneratePasswordHash(password)
	fmt.Println(hash)
	if err != nil {
		t.Fatalf("Error generating password hash: %v", err)
	}

	// Ensure that the generated hash is not empty
	if hash == "" {
		t.Error("Generated hash is empty")
	}

	// Verify the password with the generated hash
	match, err := VerfiyPassword(password, hash)
	if err != nil {
		t.Fatalf("Error verifying password: %v", err)
	}

	// Ensure that the verification succeeds
	if !match {
		t.Error("Password verification failed for a valid password")
	}

	// Test case 2: Invalid password verification
	invalidPassword := "wrongPassword456"
	match, err = VerfiyPassword(invalidPassword, hash)
	if err != nil {
		t.Fatalf("Error verifying password: %v", err)
	}

	// Ensure that the verification fails for an invalid password
	if match {
		t.Error("Password verification succeeded for an invalid password")
	}

	// Test case 3: Invalid hash format
	invalidHash := "invalidHashFormat"
	match, err = VerfiyPassword(password, invalidHash)
	if err == nil || match {
		t.Error("Password verification succeeded for an invalid hash format")
	}
}

func TestDecodeHash(t *testing.T) {
	// Test the decodeHash function

	// Test case 1: Valid encoded hash
	validEncodedHash := "$argon2id$v=19$m=65536,t=1,p=4$vZCUJbj5S9vH2WXjMet1LQ$sGYTyjW23Q5XmifAsqHxlzNxyP2GRKJ791quc/5nYHg"
	p, salt, hash, err := decodeHash(validEncodedHash)
	if err != nil {
		t.Fatalf("Error decoding valid hash: %v", err)
	}

	// Ensure that the decoded parameters are correct
	if p.memory != 65536 || p.time != 1 || p.parallelism != 4 || p.saltLength != 16 || p.keyLength != 32 {
		t.Error("Decoded parameters are incorrect")
	}

	// Ensure that the decoded salt and hash are not empty
	if len(salt) == 0 || len(hash) == 0 {
		t.Error("Decoded salt or hash is empty")
	}

	// Test case 2: Invalid encoded hash
	invalidEncodedHash := "invalidEncodedHash"
	_, _, _, err = decodeHash(invalidEncodedHash)
	if err == nil {
		t.Error("Decoding should fail for an invalid encoded hash")
	}
}
