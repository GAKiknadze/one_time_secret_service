package storage

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Record struct {
	ID        uint32    `gorm:"primaryKey"`
	Data      []byte    `gorm:"not null"`
	CreatedAt time.Time `gorm:"index"`
}

var ErrNotFound = errors.New("record not found")

type StorageDatabase struct {
	db *gorm.DB
}

// NewStorageDatabase инициализирует хранилище и создаёт таблицу при необходимости.
func NewStorageDatabase(db *gorm.DB) (*StorageDatabase, error) {
	if err := db.AutoMigrate(&Record{}); err != nil {
		return nil, err
	}
	return &StorageDatabase{db: db}, nil
}

// Save сохраняет или обновляет запись по ID.
func (s *StorageDatabase) Save(id uint32, data []byte) error {
	record := Record{
		ID:   id,
		Data: data,
	}
	result := s.db.Save(&record)
	return result.Error
}

// Get возвращает данные по ID.
func (s *StorageDatabase) Get(id uint32) ([]byte, error) {
	var record Record
	result := s.db.First(&record, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return record.Data, nil
}

// DeleteById удаляет запись по ID.
func (s *StorageDatabase) DeleteById(id uint32) error {
	result := s.db.Delete(&Record{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete удаляет все записи, созданные раньше, чем (сейчас - duration).
func (s *StorageDatabase) Delete(duration time.Duration) error {
	cutoff := time.Now().Add(-duration)
	result := s.db.Where("created_at < ?", cutoff).Delete(&Record{})
	return result.Error
}
