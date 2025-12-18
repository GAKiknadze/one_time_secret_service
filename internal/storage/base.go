package storage

import "time"

type Storage interface {
	Save(id uint32, data []byte) error
	Get(id uint32) ([]byte, error)
	DeleteById(id uint32) error
	Delete(duration time.Duration) error
}
