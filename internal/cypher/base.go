package cypher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

type Cypher interface {
	Encrypt(key, data string) []byte
	Decrypt(key string, data []byte) (string, error)
	GenerateKey(bits int) (string, error)
}

type CypherBase struct{}

// Encrypt шифрует данные с использованием AES-GCM
// Ключ должен быть валидной hex-строкой длиной 32, 48 или 64 символа (128, 192 или 256 бит)
func (c *CypherBase) Encrypt(key, data string) []byte {
	// Декодируем hex-ключ в байты
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		// В реальном приложении лучше возвращать ошибку, но сигнатура метода этого не позволяет
		// Поэтому возвращаем nil или паникуем
		panic("invalid hex key format")
	}

	// Проверяем длину ключа
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		panic("key length must be 16, 24, or 32 bytes (32, 48, or 64 hex chars)")
	}

	// Создаем AES-шифр
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		panic(err)
	}

	// Создаем GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	// Генерируем nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	// Шифруем данные
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)
	return ciphertext
}

// Decrypt расшифровывает данные, зашифрованные с помощью AES-GCM
// Ключ должен быть валидной hex-строкой длиной 32, 48 или 64 символа
func (c *CypherBase) Decrypt(key string, data []byte) (string, error) {
	// Декодируем hex-ключ в байты
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return "", errors.New("invalid hex key format")
	}

	// Проверяем длину ключа
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", errors.New("key length must be 16, 24, or 32 bytes")
	}

	// Создаем AES-шифр
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Создаем GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Извлекаем nonce из начала данных
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

// GenerateKey генерирует случайный криптографический ключ заданной битовой длины.
// Поддерживаются длины: 128, 192 и 256 бит.
// Возвращает ключ в виде hex-строки (длина строки: 32, 48 или 64 символа).
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
