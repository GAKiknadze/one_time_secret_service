package storage

import "time"

// Storage defines the interface for secret data persistence.
type Storage interface {
	// Save stores encrypted secret data with the given ID
	Save(id uint32, data []byte) error

	// Get retrieves encrypted secret data by ID
	Get(id uint32) ([]byte, error)

	// DeleteById removes a secret by its ID
	DeleteById(id uint32) error

	// Delete removes all secrets created before (now - duration)
	Delete(duration time.Duration) error
}
