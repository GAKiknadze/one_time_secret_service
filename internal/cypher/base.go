package cypher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// Cypher defines the interface for encryption and decryption operations.
type Cypher interface {
	// Encrypt encrypts plaintext using the provided hex-encoded key and returns ciphertext
	Encrypt(key, data string) []byte

	// Decrypt decrypts ciphertext using the provided hex-encoded key and returns plaintext
	Decrypt(key string, data []byte) (string, error)

	// GenerateKey generates a random cryptographic key of the specified bit length
	GenerateKey(bits int) (string, error)
}

// CypherBase provides AES-GCM encryption/decryption implementation
type CypherBase struct{}

// Encrypt encrypts data using AES-GCM.
// Key must be a valid hex string of length 32, 48, or 64 characters (128, 192, or 256 bits).
func (c *CypherBase) Encrypt(key, data string) []byte {
	// Decode hex key to bytes
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		// In a real application, it's better to return an error, but the method signature doesn't allow it.
		// So we return nil or panic.
		panic("invalid hex key format")
	}

	// Validate key length
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		panic("key length must be 16, 24, or 32 bytes (32, 48, or 64 hex chars)")
	}

	// Create AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	// Generate nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return ciphertext
}

// Decrypt decrypts data encrypted with AES-GCM.
// Key must be a valid hex string of length 32, 48, or 64 characters.
func (c *CypherBase) Decrypt(key string, data []byte) (string, error) {
	// Decode hex key to bytes
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return "", errors.New("invalid hex key format")
	}

	// Validate key length
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", errors.New("key length must be 16, 24, or 32 bytes")
	}

	// Create AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract nonce from the beginning of data
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateKey generates a random cryptographic key of the specified bit length.
// Supported lengths: 128, 192, and 256 bits.
// Returns the key as a hex string (string length: 32, 48, or 64 characters).
func (c *CypherBase) GenerateKey(bits int) (string, error) {
	var byteLen int
	switch bits {
	case 128:
		byteLen = 16
	case 192:
		byteLen = 24
	case 256:
		byteLen = 32
	default:
		return "", errors.New("unsupported key size: use 128, 192, or 256 bits")
	}

	key := make([]byte, byteLen)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}
