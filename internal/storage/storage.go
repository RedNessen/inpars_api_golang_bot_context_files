package storage

import (
	"time"

	"github.com/RedNessen/inpars-telegram-bot/internal/inpars"
)

// Storage интерфейс для работы с хранилищем объявлений
type Storage interface {
	// GetEstate получить snapshot объявления по ID
	GetEstate(id int) (*EstateSnapshot, error)

	// SaveEstate сохранить новое объявление
	SaveEstate(estate *inpars.Estate) error

	// UpdateEstate обновить существующее объявление
	UpdateEstate(snapshot *EstateSnapshot) error

	// UpdateLastSeen обновить время последнего обнаружения
	UpdateLastSeen(id int) error

	// LogChange записать изменение в историю
	LogChange(estateID int, field, oldVal, newVal string) error

	// GetChangeHistory получить историю изменений объявления
	GetChangeHistory(estateID int, limit int) ([]EstateChange, error)

	// CleanupOldEstates удалить объявления старше указанной даты
	CleanupOldEstates(olderThan time.Time) (int, error)

	// GetStats получить статистику БД
	GetStats() (*Stats, error)

	// Close закрыть соединение
	Close() error
}

// Stats статистика базы данных
type Stats struct {
	TotalEstates    int
	NewToday        int
	UpdatedToday    int
	DatabaseSizeMB  float64
	LastCleanupTime *time.Time
}
