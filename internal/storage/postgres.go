package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/RedNessen/inpars-telegram-bot/internal/inpars"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// PostgresStorage реализация Storage для PostgreSQL
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage создает новое хранилище PostgreSQL
func NewPostgresStorage(databaseURL string) (*PostgresStorage, error) {
	// Подключение к БД
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL")

	storage := &PostgresStorage{db: db}

	// Запуск миграций
	if err := storage.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed")

	return storage, nil
}

// runMigrations выполняет миграции базы данных
func (s *PostgresStorage) runMigrations() error {
	// Создаем источник миграций из embed.FS
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	// Создаем драйвер для PostgreSQL
	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Создаем мигратор
	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Выполняем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

// GetEstate получить snapshot объявления по ID
func (s *PostgresStorage) GetEstate(id int) (*EstateSnapshot, error) {
	query := `
		SELECT id, cost, floor, floors, sq, image_count, title, address, url, phones,
		       region_id, city_id, type_ad, section_id, category_id, agent,
		       first_seen_at, last_seen_at, last_sent_at, api_updated_at,
		       created_at, updated_at
		FROM estates
		WHERE id = $1
	`

	var snapshot EstateSnapshot
	err := s.db.QueryRow(query, id).Scan(
		&snapshot.ID, &snapshot.Cost, &snapshot.Floor, &snapshot.Floors,
		&snapshot.Sq, &snapshot.ImageCount, &snapshot.Title, &snapshot.Address,
		&snapshot.URL, &snapshot.Phones, &snapshot.RegionID, &snapshot.CityID,
		&snapshot.TypeAd, &snapshot.SectionID, &snapshot.CategoryID, &snapshot.Agent,
		&snapshot.FirstSeenAt, &snapshot.LastSeenAt, &snapshot.LastSentAt,
		&snapshot.APIUpdatedAt, &snapshot.CreatedAt, &snapshot.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Не найдено
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get estate: %w", err)
	}

	return &snapshot, nil
}

// SaveEstate сохранить новое объявление
func (s *PostgresStorage) SaveEstate(estate *inpars.Estate) error {
	snapshot := FromEstate(estate)

	query := `
		INSERT INTO estates (
			id, cost, floor, floors, sq, image_count, title, address, url, phones,
			region_id, city_id, type_ad, section_id, category_id, agent,
			first_seen_at, last_seen_at, api_updated_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21
		)
	`

	_, err := s.db.Exec(query,
		snapshot.ID, snapshot.Cost, snapshot.Floor, snapshot.Floors,
		snapshot.Sq, snapshot.ImageCount, snapshot.Title, snapshot.Address,
		snapshot.URL, snapshot.Phones, snapshot.RegionID, snapshot.CityID,
		snapshot.TypeAd, snapshot.SectionID, snapshot.CategoryID, snapshot.Agent,
		snapshot.FirstSeenAt, snapshot.LastSeenAt, snapshot.APIUpdatedAt,
		snapshot.CreatedAt, snapshot.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save estate: %w", err)
	}

	return nil
}

// UpdateEstate обновить существующее объявление
func (s *PostgresStorage) UpdateEstate(snapshot *EstateSnapshot) error {
	query := `
		UPDATE estates SET
			cost = $2,
			floor = $3,
			floors = $4,
			sq = $5,
			image_count = $6,
			title = $7,
			address = $8,
			url = $9,
			phones = $10,
			last_seen_at = $11,
			last_sent_at = $12,
			api_updated_at = $13,
			updated_at = $14
		WHERE id = $1
	`

	_, err := s.db.Exec(query,
		snapshot.ID, snapshot.Cost, snapshot.Floor, snapshot.Floors,
		snapshot.Sq, snapshot.ImageCount, snapshot.Title, snapshot.Address,
		snapshot.URL, snapshot.Phones, snapshot.LastSeenAt, snapshot.LastSentAt,
		snapshot.APIUpdatedAt, time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to update estate: %w", err)
	}

	return nil
}

// UpdateLastSeen обновить время последнего обнаружения
func (s *PostgresStorage) UpdateLastSeen(id int) error {
	query := `UPDATE estates SET last_seen_at = $1 WHERE id = $2`
	_, err := s.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update last_seen_at: %w", err)
	}
	return nil
}

// LogChange записать изменение в историю
func (s *PostgresStorage) LogChange(estateID int, field, oldVal, newVal string) error {
	query := `
		INSERT INTO estate_changes (estate_id, field_name, old_value, new_value, changed_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.Exec(query, estateID, field, oldVal, newVal, time.Now())
	if err != nil {
		return fmt.Errorf("failed to log change: %w", err)
	}

	return nil
}

// GetChangeHistory получить историю изменений объявления
func (s *PostgresStorage) GetChangeHistory(estateID int, limit int) ([]EstateChange, error) {
	query := `
		SELECT id, estate_id, field_name, old_value, new_value, changed_at
		FROM estate_changes
		WHERE estate_id = $1
		ORDER BY changed_at DESC
		LIMIT $2
	`

	rows, err := s.db.Query(query, estateID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get change history: %w", err)
	}
	defer rows.Close()

	var changes []EstateChange
	for rows.Next() {
		var change EstateChange
		err := rows.Scan(
			&change.ID, &change.EstateID, &change.FieldName,
			&change.OldValue, &change.NewValue, &change.ChangedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan change: %w", err)
		}
		changes = append(changes, change)
	}

	return changes, nil
}

// CleanupOldEstates удалить объявления старше указанной даты
func (s *PostgresStorage) CleanupOldEstates(olderThan time.Time) (int, error) {
	query := `DELETE FROM estates WHERE last_seen_at < $1`

	result, err := s.db.Exec(query, olderThan)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old estates: %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return int(count), nil
}

// GetStats получить статистику БД
func (s *PostgresStorage) GetStats() (*Stats, error) {
	stats := &Stats{}

	// Общее количество
	err := s.db.QueryRow("SELECT COUNT(*) FROM estates").Scan(&stats.TotalEstates)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Новых за сегодня
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM estates
		WHERE first_seen_at >= CURRENT_DATE
	`).Scan(&stats.NewToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get new today count: %w", err)
	}

	// Обновленных за сегодня
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM estates
		WHERE last_sent_at >= CURRENT_DATE AND first_seen_at < CURRENT_DATE
	`).Scan(&stats.UpdatedToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated today count: %w", err)
	}

	// Размер БД (приблизительно)
	var sizeBytes int64
	err = s.db.QueryRow(`
		SELECT pg_database_size(current_database())
	`).Scan(&sizeBytes)
	if err == nil {
		stats.DatabaseSizeMB = float64(sizeBytes) / 1024.0 / 1024.0
	}

	return stats, nil
}

// Close закрыть соединение
func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
