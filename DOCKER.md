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
  inpars-bot:
    image: ghcr.io/rednessen/inpars_api_golang_bot_context_files:latest
    container_name: inpars-telegram-bot
    restart: unless-stopped
    environment:
      # Обязательные параметры
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}

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

      # Часовой пояс
      - TZ=Europe/Kaliningrad

    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Использование с .env файлом

Создайте `.env` файл в директории с docker-compose.yml:

```bash
# .env
TELEGRAM_BOT_TOKEN=your_bot_token_here
INPARS_API_TOKEN=aEcS9UfAagInparSiv23aoa_vPzxqWvm

DEFAULT_REGIONS=39
MIN_COST=25000
MAX_COST=50000
FLOOR_MIN=3
FLOOR_MAX=0

TYPE_AD=1
SELLER_TYPES=1,2,3
POLL_INTERVAL=60
MAX_LISTINGS=50
```

Упрощенный docker-compose.yml:

```yaml
version: '3.8'

services:
  inpars-bot:
    image: ghcr.io/rednessen/inpars_api_golang_bot_context_files:latest
    container_name: inpars-telegram-bot
    restart: unless-stopped
    env_file:
      - .env
    environment:
      - TZ=Europe/Kaliningrad
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
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

## Проверка работы

```bash
# Проверить что контейнер запущен
docker ps | grep inpars

# Посмотреть логи запуска
docker logs inpars-telegram-bot

# Успешный запуск выглядит так:
# 2026/01/15 19:37:19 Starting InPars Telegram Bot...
# 2026/01/15 19:37:19 Configuration loaded successfully
# 2026/01/15 19:37:20 Telegram bot authorized as @YourBot
# 2026/01/15 19:37:20 Monitor initialized
# 2026/01/15 19:37:20 Bot and monitor are running.
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
  - DEFAULT_REGIONS=39
  - TYPE_AD=1
  - MIN_COST=25000
  - MAX_COST=50000
  - FLOOR_MIN=3
```

### Мониторинг покупки в Москве

```yaml
environment:
  - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
  - DEFAULT_REGIONS=77
  - TYPE_AD=2
  - MIN_COST=5000000
  - MAX_COST=15000000
```

### Несколько регионов

```yaml
environment:
  - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
  - DEFAULT_REGIONS=39,77,78  # Калининград, Москва, СПб
  - TYPE_AD=1
  - SELLER_TYPES=1  # Только собственники
```

## Troubleshooting

### Ошибка: TELEGRAM_BOT_TOKEN is required

Убедитесь что переменная `TELEGRAM_BOT_TOKEN` установлена в `.env` или `environment`.

### Бот не присылает объявления

1. Убедитесь что вы отправили `/start` боту в Telegram
2. Проверьте логи: `docker logs inpars-telegram-bot`
3. Убедитесь что фильтры не слишком жесткие (есть ли объявления под эти критерии)

### Rate limit exceeded

Увеличьте `POLL_INTERVAL` до 120 секунд или больше.

## Безопасность

- ✅ Образ использует непривилегированного пользователя (UID 1000)
- ✅ Minimal Alpine-based image
- ✅ Multi-stage build для минимального размера
- ✅ Поддержка amd64 и arm64 архитектур
