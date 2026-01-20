# Docker Deployment Guide

## Использование образа из GHCR

Этот бот автоматически публикуется в GitHub Container Registry (GHCR) при каждом push в main или при создании тега.

### Доступные образы

```
ghcr.io/rednessen/inpars_api_golang_bot_context_files:latest
ghcr.io/rednessen/inpars_api_golang_bot_context_files:main
ghcr.io/rednessen/inpars_api_golang_bot_context_files:v1.0.0  # для тегированных версий
```

## Пример использования в docker-compose

### В вашем репозитории инфраструктуры создайте docker-compose.yml:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: inpars-postgres
    restart: unless-stopped
    environment:
      - POSTGRES_USER=inpars
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-changeme}
      - POSTGRES_DB=inpars
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U inpars"]
      interval: 10s
      timeout: 5s
      retries: 5

  inpars-bot:
    image: ghcr.io/rednessen/inpars_api_golang_bot_context_files:latest
    container_name: inpars-telegram-bot
    restart: unless-stopped
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      # Обязательные параметры
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
      - DATABASE_URL=postgres://inpars:${POSTGRES_PASSWORD:-changeme}@postgres:5432/inpars?sslmode=disable

      # InPars API (опционально, есть тестовый токен по умолчанию)
      - INPARS_API_TOKEN=${INPARS_API_TOKEN:-aEcS9UfAagInparSiv23aoa_vPzxqWvm}

      # Настройки мониторинга
      - POLL_INTERVAL=60
      - MAX_LISTINGS=50

      # Фильтры региона
      - DEFAULT_REGIONS=39  # Калининградская область
      - DEFAULT_CITIES=

      # Тип объявлений
      - TYPE_AD=1  # 1=сдам, 2=продам, 3=сниму, 4=куплю
      - SELLER_TYPES=1,2,3  # 1=собственник, 2=агент, 3=застройщик

      # Ценовые фильтры
      - MIN_COST=25000
      - MAX_COST=50000

      # Фильтры этажей
      - FLOOR_MIN=3
      - FLOOR_MAX=0  # 0 = без ограничения

      # Настройки очистки
      - CLEANUP_DAYS=14  # Удалять объявления старше 14 дней
      - CLEANUP_INTERVAL=24  # Запускать очистку каждые 24 часа

      # Часовой пояс
      - TZ=Europe/Kaliningrad

    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

volumes:
  postgres_data:
```

### Использование с .env файлом

Создайте `.env` файл в директории с docker-compose.yml:

```bash
# .env
TELEGRAM_BOT_TOKEN=your_bot_token_here
INPARS_API_TOKEN=aEcS9UfAagInparSiv23aoa_vPzxqWvm

# Database
POSTGRES_PASSWORD=changeme

# Filters
DEFAULT_REGIONS=39
MIN_COST=25000
MAX_COST=50000
FLOOR_MIN=3
FLOOR_MAX=0

# Settings
TYPE_AD=1
SELLER_TYPES=1,2,3
POLL_INTERVAL=60
MAX_LISTINGS=50

# Cleanup
CLEANUP_DAYS=14
CLEANUP_INTERVAL=24
```

Упрощенный docker-compose.yml:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: inpars-postgres
    restart: unless-stopped
    environment:
      - POSTGRES_USER=inpars
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-changeme}
      - POSTGRES_DB=inpars
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U inpars"]
      interval: 10s
      timeout: 5s
      retries: 5

  inpars-bot:
    image: ghcr.io/rednessen/inpars_api_golang_bot_context_files:latest
    container_name: inpars-telegram-bot
    restart: unless-stopped
    depends_on:
      postgres:
        condition: service_healthy
    env_file:
      - .env
    environment:
      - DATABASE_URL=postgres://inpars:${POSTGRES_PASSWORD:-changeme}@postgres:5432/inpars?sslmode=disable
      - TZ=Europe/Kaliningrad
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

volumes:
  postgres_data:
```

## Запуск

```bash
# Скачать образ
docker pull ghcr.io/rednessen/inpars_api_golang_bot_context_files:latest

# Запустить через docker-compose
docker-compose up -d

# Посмотреть логи
docker-compose logs -f inpars-bot

# Остановить
docker-compose down
```

## Переменные окружения

Все доступные переменные:

| Переменная | Обязательная | Описание | По умолчанию |
|------------|--------------|----------|--------------|
| `TELEGRAM_BOT_TOKEN` | ✅ Да | Токен Telegram бота от @BotFather | - |
| `DATABASE_URL` | ✅ Да | Строка подключения к PostgreSQL | - |
| `INPARS_API_TOKEN` | ❌ Нет | Токен InPars API | Тестовый токен |
| `POLL_INTERVAL` | ❌ Нет | Интервал проверки (сек) | 60 |
| `MAX_LISTINGS` | ❌ Нет | Макс. объявлений за запрос | 50 |
| `DEFAULT_REGIONS` | ❌ Нет | ID регионов (через запятую) | 77 (Москва) |
| `DEFAULT_CITIES` | ❌ Нет | ID городов (через запятую) | - |
| `TYPE_AD` | ❌ Нет | Типы объявлений | 1 (аренда) |
| `SELLER_TYPES` | ❌ Нет | Типы продавцов | 1,2,3 (все) |
| `MIN_COST` | ❌ Нет | Минимальная цена | 0 |
| `MAX_COST` | ❌ Нет | Максимальная цена | 0 |
| `FLOOR_MIN` | ❌ Нет | Минимальный этаж | 0 |
| `FLOOR_MAX` | ❌ Нет | Максимальный этаж | 0 |
| `CLEANUP_DAYS` | ❌ Нет | Удалять объявления старше N дней | 14 |
| `CLEANUP_INTERVAL` | ❌ Нет | Интервал очистки (часы) | 24 |

## Проверка работы

```bash
# Проверить что контейнеры запущены
docker ps | grep inpars

# Посмотреть логи базы данных
docker logs inpars-postgres

# Посмотреть логи бота
docker logs inpars-telegram-bot

# Успешный запуск выглядит так:
# 2026/01/20 19:37:19 Starting InPars Telegram Bot...
# 2026/01/20 19:37:19 Configuration loaded successfully
# 2026/01/20 19:37:20 Database connection established
# 2026/01/20 19:37:20 Running migrations...
# 2026/01/20 19:37:20 Migrations completed successfully
# 2026/01/20 19:37:20 InPars API client initialized
# 2026/01/20 19:37:20 Telegram bot authorized as @YourBot
# 2026/01/20 19:37:20 Monitor initialized
# 2026/01/20 19:37:20 Cleanup routine started with interval: 24 hours
# 2026/01/20 19:37:20 Bot and monitor are running. Press Ctrl+C to stop.
```

## Обновление образа

```bash
# Остановить текущий контейнер
docker-compose down

# Скачать последнюю версию
docker pull ghcr.io/rednessen/inpars_api_golang_bot_context_files:latest

# Запустить заново
docker-compose up -d
```

## Примеры конфигураций

### Мониторинг аренды в Калининграде

```yaml
environment:
  - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
  - DATABASE_URL=postgres://inpars:${POSTGRES_PASSWORD}@postgres:5432/inpars?sslmode=disable
  - DEFAULT_REGIONS=39
  - TYPE_AD=1
  - MIN_COST=25000
  - MAX_COST=50000
  - FLOOR_MIN=3
  - CLEANUP_DAYS=14
```

### Мониторинг покупки в Москве

```yaml
environment:
  - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
  - DATABASE_URL=postgres://inpars:${POSTGRES_PASSWORD}@postgres:5432/inpars?sslmode=disable
  - DEFAULT_REGIONS=77
  - TYPE_AD=2
  - MIN_COST=5000000
  - MAX_COST=15000000
  - CLEANUP_DAYS=30  # Для продажи храним дольше
```

### Несколько регионов

```yaml
environment:
  - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
  - DATABASE_URL=postgres://inpars:${POSTGRES_PASSWORD}@postgres:5432/inpars?sslmode=disable
  - DEFAULT_REGIONS=39,77,78  # Калининград, Москва, СПб
  - TYPE_AD=1
  - SELLER_TYPES=1  # Только собственники
  - CLEANUP_DAYS=7  # Более частая очистка при множестве регионов
```

## Troubleshooting

### Ошибка: TELEGRAM_BOT_TOKEN is required

Убедитесь что переменная `TELEGRAM_BOT_TOKEN` установлена в `.env` или `environment`.

### Ошибка: DATABASE_URL is required

Убедитесь что переменная `DATABASE_URL` указана в конфигурации. Для docker-compose она должна ссылаться на сервис postgres:

```yaml
DATABASE_URL=postgres://inpars:${POSTGRES_PASSWORD}@postgres:5432/inpars?sslmode=disable
```

### Ошибка подключения к базе данных

1. Проверьте что контейнер PostgreSQL запущен: `docker ps | grep postgres`
2. Проверьте логи PostgreSQL: `docker logs inpars-postgres`
3. Убедитесь что бот ждет готовности БД через `depends_on` с `condition: service_healthy`

### Бот не присылает объявления

1. Убедитесь что вы отправили `/start` боту в Telegram
2. Проверьте логи: `docker logs inpars-telegram-bot`
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
