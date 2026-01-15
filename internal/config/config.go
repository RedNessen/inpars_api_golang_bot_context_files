package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config содержит конфигурацию приложения
type Config struct {
	// Telegram Bot
	TelegramToken string

	// InPars API
	InParsToken string

	// Настройки мониторинга
	PollInterval    int   // Интервал опроса API в секундах
	MaxListings     int   // Максимальное количество объявлений за один запрос

	// Фильтры по умолчанию
	DefaultRegions  []int // ID регионов для мониторинга
	DefaultCities   []int // ID городов
	TypeAd          []int // Типы объявлений (1-сдам по умолчанию)
	SellerTypes     []int // Типы продавцов (1,2,3 - все)

	// Лимиты
	MinCost         int
	MaxCost         int
}

// LoadFromEnv загружает конфигурацию из переменных окружения
func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		TelegramToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		InParsToken:     getEnvOrDefault("INPARS_API_TOKEN", "aEcS9UfAagInparSiv23aoa_vPzxqWvm"), // Тестовый токен по умолчанию
		PollInterval:    getEnvAsInt("POLL_INTERVAL", 60),    // 60 секунд по умолчанию
		MaxListings:     getEnvAsInt("MAX_LISTINGS", 50),     // 50 объявлений (лимит для тестового токена)
		DefaultRegions:  getEnvAsIntSlice("DEFAULT_REGIONS", []int{77}), // Москва по умолчанию
		DefaultCities:   getEnvAsIntSlice("DEFAULT_CITIES", []int{}),
		TypeAd:          getEnvAsIntSlice("TYPE_AD", []int{1}), // 1 - сдам (аренда)
		SellerTypes:     getEnvAsIntSlice("SELLER_TYPES", []int{1, 2, 3}), // Все типы
		MinCost:         getEnvAsInt("MIN_COST", 0),
		MaxCost:         getEnvAsInt("MAX_COST", 0),
	}

	// Валидация обязательных полей
	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	if cfg.InParsToken == "" {
		return nil, fmt.Errorf("INPARS_API_TOKEN is required")
	}

	return cfg, nil
}

// getEnvOrDefault возвращает значение переменной окружения или значение по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt возвращает значение переменной окружения как int
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnvAsIntSlice возвращает значение переменной окружения как []int
// Формат: "1,2,3" или "77"
func getEnvAsIntSlice(key string, defaultValue []int) []int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	parts := strings.Split(valueStr, ",")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		value, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		result = append(result, value)
	}

	if len(result) == 0 {
		return defaultValue
	}
	return result
}
