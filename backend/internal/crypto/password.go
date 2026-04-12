package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/argon2"
)

// Argon2id parameters tuned to OWASP recommended minimums.
// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
const (
	argon2Time    uint32 = 3
	argon2Memory  uint32 = 64 * 1024 // 64 MB
	argon2Threads uint8  = 2
	argon2KeyLen  uint32 = 32
	saltByteLen          = 16
)

// GenerateSalt creates a cryptographically random 16-byte salt and returns it
// as a lowercase hex-encoded string.
func GenerateSalt() (string, error) {
	b := make([]byte, saltByteLen)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("crypto: failed to generate random salt")
	}
	return hex.EncodeToString(b), nil
}

// HashPassword derives an Argon2id key from the plain-text password and a
// hex-encoded salt, returning the result as a hex string.
//
//	argon2id(password, salt, time=3, mem=64MB, p=2, keyLen=32)
func HashPassword(password, saltHex string) (string, error) {
	saltBytes, err := hex.DecodeString(saltHex)
	if err != nil {
		return "", errors.New("crypto: invalid salt encoding")
	}
	hash := argon2.IDKey(
		[]byte(password),
		saltBytes,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)
	return hex.EncodeToString(hash), nil
}

// VerifyPassword checks whether the plain-text password matches the stored
// Argon2id hash using constant-time comparison to prevent timing attacks.
// Returns true only when the password is correct.
func VerifyPassword(password, saltHex, storedHashHex string) bool {
	expectedHashHex, err := HashPassword(password, saltHex)
	if err != nil {
		return false
	}

	storedBytes, err := hex.DecodeString(storedHashHex)
	if err != nil {
		return false
	}
	expectedBytes, err := hex.DecodeString(expectedHashHex)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(storedBytes, expectedBytes) == 1
}
