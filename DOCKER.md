# Docker Deployment Guide

## Использование образа из GHCR

Этот бот автоматически публикуется в GitHub Container Registry (GHCR) при каждом push в main или при создании тега.

### Доступные образы

```
ghcr.io/rednessen/inpars_api_golang_bot_context_files:latest
ghcr.io/rednessen/inpars_api_golang_bot_context_files:main
ghcr.io/rednessen/inpars_api_golang_bot_context_files:v1.0.0  # для тегированных версий
```

## Требования для развертывания

Бот требует:
1. **PostgreSQL 14+** для персистентного хранения данных
2. Переменные окружения (см. ниже)

**Важно:** Миграции базы данных выполняются автоматически при запуске приложения.

## Переменные окружения

### Обязательные переменные

| Переменная | Описание | Пример |
|------------|----------|--------|
| `TELEGRAM_BOT_TOKEN` | Токен Telegram бота от @BotFather | `1234567890:ABC...` |
| `DATABASE_URL` | Строка подключения к PostgreSQL | `postgres://user:pass@host:5432/db?sslmode=disable` |

### Опциональные переменные

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `INPARS_API_TOKEN` | Токен InPars API | Тестовый токен |
| `POLL_INTERVAL` | Интервал проверки (сек) | 60 |
| `MAX_LISTINGS` | Макс. объявлений за запрос | 50 |
| `DEFAULT_REGIONS` | ID регионов (через запятую) | 77 (Москва) |
| `DEFAULT_CITIES` | ID городов (через запятую) | - |
| `TYPE_AD` | Типы объявлений: 1-сдам, 2-продам, 3-сниму, 4-куплю | 1 |
| `SELLER_TYPES` | Типы продавцов: 1-собственник, 2-агент, 3-застройщик | 1,2,3 |
| `MIN_COST` | Минимальная цена | 0 |
| `MAX_COST` | Максимальная цена | 0 |
| `FLOOR_MIN` | Минимальный этаж | 0 |
| `FLOOR_MAX` | Максимальный этаж | 0 |
| `CLEANUP_DAYS` | Удалять объявления старше N дней | 14 |
| `CLEANUP_INTERVAL` | Интервал очистки (часы) | 24 |

## Формат DATABASE_URL

```
postgres://username:password@hostname:port/database?sslmode=disable
```

Примеры:
- Локальный PostgreSQL: `postgres://inpars:password@localhost:5432/inpars?sslmode=disable`
- Docker Compose сервис: `postgres://inpars:password@postgres:5432/inpars?sslmode=disable`
- Удаленный сервер: `postgres://inpars:password@db.example.com:5432/inpars?sslmode=disable`

## Проверка работы

После запуска проверьте логи контейнера. Успешный запуск выглядит так:

```
2026/01/20 19:37:19 Starting InPars Telegram Bot...
2026/01/20 19:37:19 Configuration loaded successfully
2026/01/20 19:37:20 Database connection established
2026/01/20 19:37:20 Running migrations...
2026/01/20 19:37:20 Migrations completed successfully
2026/01/20 19:37:20 InPars API client initialized
2026/01/20 19:37:20 Telegram bot authorized as @YourBot
2026/01/20 19:37:20 Monitor initialized
2026/01/20 19:37:20 Cleanup routine started with interval: 24 hours
2026/01/20 19:37:20 Bot and monitor are running. Press Ctrl+C to stop.
```

## Troubleshooting

### Ошибка: TELEGRAM_BOT_TOKEN is required

Убедитесь что переменная `TELEGRAM_BOT_TOKEN` установлена в конфигурации.

### Ошибка: DATABASE_URL is required

Убедитесь что переменная `DATABASE_URL` указана в конфигурации и корректно отформатирована.

### Ошибка подключения к базе данных

1. Проверьте что PostgreSQL запущен и доступен по указанному адресу
2. Проверьте правильность credentials в DATABASE_URL
3. Проверьте что база данных создана
4. Если используется docker-compose, убедитесь что бот ждет готовности БД через `depends_on` с `condition: service_healthy`

### Бот не присылает объявления

1. Убедитесь что вы отправили `/start` боту в Telegram
2. Проверьте логи контейнера
3. Убедитесь что фильтры не слишком жесткие (есть ли объявления под эти критерии)
4. Проверьте что база данных доступна и миграции выполнены успешно

### Rate limit exceeded

Увеличьте `POLL_INTERVAL` до 120 секунд или больше.

### База данных растет слишком быстро

Уменьшите `CLEANUP_DAYS` для более частой очистки старых объявлений. Например, для активного мониторинга нескольких регионов используйте 7 дней вместо 14.

## Безопасность

- ✅ Образ использует непривилегированного пользователя (UID 1000)
- ✅ Minimal Alpine-based image
- ✅ Multi-stage build для минимального размера
- ✅ Поддержка amd64 и arm64 архитектур

## Обновление образа

Для получения последней версии:

```bash
docker pull ghcr.io/rednessen/inpars_api_golang_bot_context_files:latest
```

После обновления образа перезапустите контейнер. Данные в PostgreSQL сохранятся.
