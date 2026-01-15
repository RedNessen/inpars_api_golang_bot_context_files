package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/RedNessen/inpars-telegram-bot/internal/config"
	"github.com/RedNessen/inpars-telegram-bot/internal/inpars"
	"github.com/RedNessen/inpars-telegram-bot/internal/telegram"
)

// Monitor отслеживает новые объявления и отправляет уведомления
type Monitor struct {
	client       *inpars.Client
	bot          *telegram.Bot
	config       *config.Config
	lastUpdateID int                // ID последнего обработанного объявления
	seenIDs      map[int]bool       // Множество уже обработанных ID
	lastUpdate   time.Time          // Время последнего обновления
}

// NewMonitor создает новый монитор
func NewMonitor(client *inpars.Client, bot *telegram.Bot, cfg *config.Config) *Monitor {
	return &Monitor{
		client:     client,
		bot:        bot,
		config:     cfg,
		seenIDs:    make(map[int]bool),
		lastUpdate: time.Now(),
	}
}

// Start запускает мониторинг
func (m *Monitor) Start() error {
	log.Println("Starting monitoring service...")

	// Инициализация: получаем последние объявления чтобы не отправлять старые при старте
	if err := m.initializeLastSeen(); err != nil {
		log.Printf("Warning: failed to initialize last seen: %v", err)
	}

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

// initializeLastSeen инициализирует список уже существующих объявлений
func (m *Monitor) initializeLastSeen() error {
	log.Println("Initializing last seen listings...")

	params := m.buildParams()
	params.Limit = m.config.MaxListings
	params.SortBy = "id_desc" // Сортируем по ID в порядке убывания

	resp, err := m.client.GetEstateList(params)
	if err != nil {
		return fmt.Errorf("failed to get initial listings: %w", err)
	}

	// Сохраняем ID существующих объявлений
	for _, estate := range resp.Data {
		m.seenIDs[estate.ID] = true
		if estate.ID > m.lastUpdateID {
			m.lastUpdateID = estate.ID
		}
	}

	log.Printf("Initialized with %d existing listings. Last ID: %d", len(resp.Data), m.lastUpdateID)
	return nil
}

// checkForNewListings проверяет наличие новых объявлений
func (m *Monitor) checkForNewListings() error {
	// Проверяем, есть ли активные чаты
	if !m.bot.HasActiveChats() {
		log.Println("No active chats, skipping check...")
		return nil
	}

	log.Println("Checking for new listings...")

	params := m.buildParams()
	params.Limit = m.config.MaxListings
	params.SortBy = "id_desc" // Сортируем по ID в порядке убывания

	// Если есть последний ID, запрашиваем объявления с ID больше него
	if m.lastUpdateID > 0 {
		params.LastID = m.lastUpdateID
		params.SortBy = "id_asc" // При использовании lastId используем сортировку по возрастанию
	}

	resp, err := m.client.GetEstateList(params)
	if err != nil {
		return fmt.Errorf("failed to get estate list: %w", err)
	}

	// Обрабатываем новые объявления
	newCount := 0
	for _, estate := range resp.Data {
		// Пропускаем уже обработанные объявления
		if m.seenIDs[estate.ID] {
			continue
		}

		// Отмечаем как обработанное
		m.seenIDs[estate.ID] = true

		// Обновляем последний ID
		if estate.ID > m.lastUpdateID {
			m.lastUpdateID = estate.ID
		}

		// Отправляем уведомление
		if err := m.bot.SendEstate(&estate); err != nil {
			log.Printf("Failed to send estate %d: %v", estate.ID, err)
			continue
		}

		newCount++
		log.Printf("Sent new listing: ID=%d, Title=%s", estate.ID, estate.Title)

		// Задержка между отправками, чтобы избежать флуда
		time.Sleep(500 * time.Millisecond)
	}

	// Очищаем старые записи из seenIDs для экономии памяти
	// Храним только последние 10000 записей
	if len(m.seenIDs) > 10000 {
		m.cleanupSeenIDs()
	}

	if newCount > 0 {
		log.Printf("Found and sent %d new listings. Last ID: %d", newCount, m.lastUpdateID)
	} else {
		log.Println("No new listings found")
	}

	// Выводим информацию о rate limiting
	if resp.Meta.RateRemaining > 0 {
		log.Printf("Rate limit: %d/%d remaining, resets in %d seconds",
			resp.Meta.RateRemaining, resp.Meta.RateLimit, resp.Meta.RateReset)
	}

	m.lastUpdate = time.Now()
	return nil
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

	return params
}

// cleanupSeenIDs очищает старые записи из seenIDs
func (m *Monitor) cleanupSeenIDs() {
	log.Println("Cleaning up old seen IDs...")

	// Оставляем только ID близкие к последнему
	minID := m.lastUpdateID - 5000
	newSeenIDs := make(map[int]bool)

	for id := range m.seenIDs {
		if id >= minID {
			newSeenIDs[id] = true
		}
	}

	m.seenIDs = newSeenIDs
	log.Printf("Cleaned up seen IDs. Current count: %d", len(m.seenIDs))
}

// GetStatus возвращает статус монитора
func (m *Monitor) GetStatus() string {
	return fmt.Sprintf(
		"Monitor Status:\n"+
			"Last Update: %s\n"+
			"Last ID: %d\n"+
			"Seen IDs: %d\n"+
			"Active Chats: %d",
		m.lastUpdate.Format("2006-01-02 15:04:05"),
		m.lastUpdateID,
		len(m.seenIDs),
		len(m.bot.GetActiveChatIDs()),
	)
}
