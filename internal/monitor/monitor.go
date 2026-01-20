package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/RedNessen/inpars-telegram-bot/internal/config"
	"github.com/RedNessen/inpars-telegram-bot/internal/inpars"
	"github.com/RedNessen/inpars-telegram-bot/internal/storage"
	"github.com/RedNessen/inpars-telegram-bot/internal/telegram"
)

// Monitor отслеживает новые объявления и отправляет уведомления
type Monitor struct {
	client       *inpars.Client
	bot          *telegram.Bot
	storage      storage.Storage
	config       *config.Config
	lastUpdate   time.Time
}

// NewMonitor создает новый монитор
func NewMonitor(client *inpars.Client, bot *telegram.Bot, store storage.Storage, cfg *config.Config) *Monitor {
	return &Monitor{
		client:     client,
		bot:        bot,
		storage:    store,
		config:     cfg,
		lastUpdate: time.Now(),
	}
}

// Start запускает мониторинг
func (m *Monitor) Start() error {
	log.Println("Starting monitoring service...")

	// Запускаем горутину для периодической очистки
	go m.startCleanupRoutine()

	// Запускаем цикл мониторинга
	ticker := time.NewTicker(time.Duration(m.config.PollInterval) * time.Second)
	defer ticker.Stop()

	log.Printf("Monitoring started with interval: %d seconds", m.config.PollInterval)

	for {
		select {
		case <-ticker.C:
			if err := m.checkForNewListings(); err != nil {
				log.Printf("Error checking for new listings: %v", err)
			}
		}
	}
}

// checkForNewListings проверяет наличие новых и обновлённых объявлений
func (m *Monitor) checkForNewListings() error {
	// Проверяем, есть ли активные чаты
	if !m.bot.HasActiveChats() {
		log.Println("No active chats, skipping check...")
		return nil
	}

	log.Println("Checking for new listings...")

	params := m.buildParams()
	params.Limit = m.config.MaxListings
	params.SortBy = "updated_desc" // Сортируем по дате обновления

	resp, err := m.client.GetEstateList(params)
	if err != nil {
		return fmt.Errorf("failed to get estate list: %w", err)
	}

	// Обрабатываем объявления
	newCount := 0
	updatedCount := 0

	for _, estate := range resp.Data {
		// Проверяем есть ли объявление в БД
		existingSnapshot, err := m.storage.GetEstate(estate.ID)
		if err != nil {
			log.Printf("Error getting estate %d from storage: %v", estate.ID, err)
			continue
		}

		if existingSnapshot == nil {
			// НОВОЕ объявление
			if err := m.handleNewEstate(&estate); err != nil {
				log.Printf("Failed to handle new estate %d: %v", estate.ID, err)
				continue
			}
			newCount++
		} else {
			// СУЩЕСТВУЮЩЕЕ объявление - проверяем изменения
			updated, err := m.handleExistingEstate(&estate, existingSnapshot)
			if err != nil {
				log.Printf("Failed to handle existing estate %d: %v", estate.ID, err)
				continue
			}
			if updated {
				updatedCount++
			}
		}

		// Задержка между отправками
		if newCount+updatedCount > 0 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	if newCount > 0 || updatedCount > 0 {
		log.Printf("Found %d new and %d updated listings", newCount, updatedCount)
	} else {
		log.Println("No new or updated listings found")
	}

	// Выводим информацию о rate limiting
	if resp.Meta.RateRemaining > 0 {
		log.Printf("Rate limit: %d/%d remaining, resets in %d seconds",
			resp.Meta.RateRemaining, resp.Meta.RateLimit, resp.Meta.RateReset)
	}

	m.lastUpdate = time.Now()
	return nil
}

// handleNewEstate обрабатывает новое объявление
func (m *Monitor) handleNewEstate(estate *inpars.Estate) error {
	// Сохраняем в БД
	if err := m.storage.SaveEstate(estate); err != nil {
		return fmt.Errorf("failed to save estate: %w", err)
	}

	// Отправляем уведомление "НОВОЕ"
	if err := m.bot.SendNewEstate(estate); err != nil {
		return fmt.Errorf("failed to send new estate notification: %w", err)
	}

	// Обновляем last_sent_at
	snapshot := storage.FromEstate(estate)
	now := time.Now()
	snapshot.LastSentAt = &now
	if err := m.storage.UpdateEstate(snapshot); err != nil {
		log.Printf("Warning: failed to update last_sent_at: %v", err)
	}

	log.Printf("🆕 Sent NEW listing: ID=%d, Title=%s, Cost=%d", estate.ID, estate.Title, estate.Cost)
	return nil
}

// handleExistingEstate обрабатывает существующее объявление
func (m *Monitor) handleExistingEstate(estate *inpars.Estate, oldSnapshot *storage.EstateSnapshot) (bool, error) {
	// Создаём новый snapshot
	newSnapshot := storage.FromEstate(estate)

	// Сравниваем с предыдущим
	changes := oldSnapshot.CompareWith(newSnapshot)

	if len(changes) == 0 {
		// Изменений нет - просто обновляем last_seen_at
		if err := m.storage.UpdateLastSeen(estate.ID); err != nil {
			return false, fmt.Errorf("failed to update last_seen: %w", err)
		}
		return false, nil
	}

	// Есть изменения!
	log.Printf("🔄 Detected changes in estate %d: %d fields changed", estate.ID, len(changes))

	// Записываем изменения в историю
	for _, change := range changes {
		if err := m.storage.LogChange(estate.ID, string(change.Field), change.OldValue, change.NewValue); err != nil {
			log.Printf("Warning: failed to log change: %v", err)
		}
	}

	// Обновляем snapshot в БД
	newSnapshot.FirstSeenAt = oldSnapshot.FirstSeenAt // Сохраняем original first_seen
	newSnapshot.LastSeenAt = time.Now()
	now := time.Now()
	newSnapshot.LastSentAt = &now

	if err := m.storage.UpdateEstate(newSnapshot); err != nil {
		return false, fmt.Errorf("failed to update estate: %w", err)
	}

	// Отправляем уведомление "ОБНОВЛЕНО"
	if err := m.bot.SendUpdatedEstate(estate, changes); err != nil {
		return false, fmt.Errorf("failed to send updated estate notification: %w", err)
	}

	log.Printf("🔄 Sent UPDATED listing: ID=%d, Changes=%v", estate.ID, formatChanges(changes))
	return true, nil
}

// startCleanupRoutine запускает периодическую очистку старых объявлений
func (m *Monitor) startCleanupRoutine() {
	ticker := time.NewTicker(time.Duration(m.config.CleanupInterval) * time.Hour)
	defer ticker.Stop()

	log.Printf("Cleanup routine started with interval: %d hours, threshold: %d days",
		m.config.CleanupInterval, m.config.CleanupDays)

	// Запускаем сразу при старте
	m.runCleanup()

	for {
		select {
		case <-ticker.C:
			m.runCleanup()
		}
	}
}

// runCleanup выполняет очистку старых объявлений
func (m *Monitor) runCleanup() {
	log.Println("Running cleanup of old estates...")

	threshold := time.Now().AddDate(0, 0, -m.config.CleanupDays)
	count, err := m.storage.CleanupOldEstates(threshold)
	if err != nil {
		log.Printf("Error during cleanup: %v", err)
		return
	}

	if count > 0 {
		log.Printf("Cleanup completed: removed %d old estates (older than %d days)", count, m.config.CleanupDays)
	} else {
		log.Println("Cleanup completed: no old estates to remove")
	}
}

// buildParams создает параметры запроса на основе конфигурации
func (m *Monitor) buildParams() *inpars.EstateListParams {
	params := &inpars.EstateListParams{
		TypeAd:     m.config.TypeAd,
		SellerType: m.config.SellerTypes,
		Expand: []string{
			"region", "city", "metro", "category",
			"material", "rentTime", "rooms", "rentTerms",
		},
	}

	// Регионы
	if len(m.config.DefaultRegions) > 0 {
		params.RegionID = m.config.DefaultRegions
	}

	// Города
	if len(m.config.DefaultCities) > 0 {
		params.CityID = m.config.DefaultCities
	}

	// Фильтры по цене
	if m.config.MinCost > 0 {
		params.CostMin = m.config.MinCost
	}
	if m.config.MaxCost > 0 {
		params.CostMax = m.config.MaxCost
	}

	// Фильтры по этажам
	if m.config.FloorMin > 0 {
		params.FloorMin = m.config.FloorMin
	}
	if m.config.FloorMax > 0 {
		params.FloorMax = m.config.FloorMax
	}

	return params
}

// GetStatus возвращает статус монитора
func (m *Monitor) GetStatus() string {
	stats, err := m.storage.GetStats()
	if err != nil {
		return fmt.Sprintf("Error getting stats: %v", err)
	}

	return fmt.Sprintf(
		"Monitor Status:\n"+
			"Last Check: %s\n"+
			"Total Estates: %d\n"+
			"New Today: %d\n"+
			"Updated Today: %d\n"+
			"Database Size: %.2f MB\n"+
			"Active Chats: %d",
		m.lastUpdate.Format("2006-01-02 15:04:05"),
		stats.TotalEstates,
		stats.NewToday,
		stats.UpdatedToday,
		stats.DatabaseSizeMB,
		len(m.bot.GetActiveChatIDs()),
	)
}

// formatChanges форматирует список изменений для лога
func formatChanges(changes []storage.Change) string {
	if len(changes) == 0 {
		return "none"
	}
	result := ""
	for i, ch := range changes {
		if i > 0 {
			result += ", "
		}
		result += string(ch.Field)
	}
	return result
}
