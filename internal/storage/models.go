package storage

import (
	"time"

	"github.com/RedNessen/inpars-telegram-bot/internal/inpars"
)

// EstateSnapshot представляет snapshot объявления в БД
type EstateSnapshot struct {
	ID         int       `db:"id"`
	Cost       int       `db:"cost"`
	Floor      int       `db:"floor"`
	Floors     int       `db:"floors"`
	Sq         float64   `db:"sq"`
	ImageCount int       `db:"image_count"`
	Title      string    `db:"title"`
	Address    string    `db:"address"`
	URL        string    `db:"url"`
	Phones     string    `db:"phones"` // JSON string

	RegionID   int `db:"region_id"`
	CityID     int `db:"city_id"`
	TypeAd     int `db:"type_ad"`
	SectionID  int `db:"section_id"`
	CategoryID int `db:"category_id"`
	Agent      int `db:"agent"`

	FirstSeenAt  time.Time  `db:"first_seen_at"`
	LastSeenAt   time.Time  `db:"last_seen_at"`
	LastSentAt   *time.Time `db:"last_sent_at"`
	APIUpdatedAt string     `db:"api_updated_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// EstateChange представляет изменение в объявлении
type EstateChange struct {
	ID        int       `db:"id"`
	EstateID  int       `db:"estate_id"`
	FieldName string    `db:"field_name"`
	OldValue  string    `db:"old_value"`
	NewValue  string    `db:"new_value"`
	ChangedAt time.Time `db:"changed_at"`
}

// ChangeType тип изменения
type ChangeType string

const (
	ChangeTypeCost       ChangeType = "cost"
	ChangeTypeFloor      ChangeType = "floor"
	ChangeTypeImageCount ChangeType = "image_count"
	ChangeTypeSq         ChangeType = "sq"
)

// Change представляет изменение для отображения пользователю
type Change struct {
	Field    ChangeType
	OldValue string
	NewValue string
}

// FromEstate создает EstateSnapshot из inpars.Estate
func FromEstate(estate *inpars.Estate) *EstateSnapshot {
	now := time.Now()

	// Преобразуем phones в JSON string
	phonesJSON := "[]"
	if len(estate.Phones) > 0 {
		phonesJSON = phonesToJSON(estate.Phones)
	}

	return &EstateSnapshot{
		ID:           estate.ID,
		Cost:         estate.Cost,
		Floor:        estate.Floor,
		Floors:       estate.Floors,
		Sq:           estate.Sq,
		ImageCount:   len(estate.Images),
		Title:        estate.Title,
		Address:      estate.Address,
		URL:          estate.URL,
		Phones:       phonesJSON,
		RegionID:     estate.RegionID,
		CityID:       estate.CityID,
		TypeAd:       estate.TypeAd,
		SectionID:    estate.SectionID,
		CategoryID:   estate.CategoryID,
		Agent:        estate.Agent,
		FirstSeenAt:  now,
		LastSeenAt:   now,
		APIUpdatedAt: estate.Updated,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// CompareWith сравнивает два snapshot и возвращает список изменений
func (s *EstateSnapshot) CompareWith(new *EstateSnapshot) []Change {
	var changes []Change

	// Проверяем цену
	if s.Cost != new.Cost {
		changes = append(changes, Change{
			Field:    ChangeTypeCost,
			OldValue: formatInt(s.Cost),
			NewValue: formatInt(new.Cost),
		})
	}

	// Проверяем этаж
	if s.Floor != new.Floor && new.Floor > 0 {
		changes = append(changes, Change{
			Field:    ChangeTypeFloor,
			OldValue: formatInt(s.Floor),
			NewValue: formatInt(new.Floor),
		})
	}

	// Проверяем количество фото
	if s.ImageCount != new.ImageCount {
		changes = append(changes, Change{
			Field:    ChangeTypeImageCount,
			OldValue: formatInt(s.ImageCount),
			NewValue: formatInt(new.ImageCount),
		})
	}

	// Проверяем площадь
	if s.Sq != new.Sq && new.Sq > 0 {
		changes = append(changes, Change{
			Field:    ChangeTypeSq,
			OldValue: formatFloat(s.Sq),
			NewValue: formatFloat(new.Sq),
		})
	}

	return changes
}

// Вспомогательные функции для форматирования
func formatInt(val int) string {
	return intToString(val)
}

func formatFloat(val float64) string {
	// Простое преобразование float64 в строку
	s := ""
	i := int(val)
	f := val - float64(i)

	s = intToString(i)
	if f > 0.001 {
		s += "."
		frac := int(f * 10)
		s += intToString(frac)
	}
	return s
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + intToString(-n)
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+(n%10))) + result
		n /= 10
	}
	return result
}

func phonesToJSON(phones []int64) string {
	if len(phones) == 0 {
		return "[]"
	}
	result := "["
	for i, phone := range phones {
		if i > 0 {
			result += ","
		}
		result += int64ToString(phone)
	}
	result += "]"
	return result
}

func int64ToString(n int64) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + int64ToString(-n)
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+(n%10))) + result
		n /= 10
	}
	return result
}
