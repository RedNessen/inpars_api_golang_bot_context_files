-- Таблица объявлений
CREATE TABLE IF NOT EXISTS estates (
    id BIGINT PRIMARY KEY,

    -- Основные данные
    cost INTEGER NOT NULL,
    floor INTEGER DEFAULT 0,
    floors INTEGER DEFAULT 0,
    sq REAL DEFAULT 0,
    image_count INTEGER DEFAULT 0,

    -- Текстовые поля
    title VARCHAR(500) DEFAULT '',
    address VARCHAR(500) DEFAULT '',
    url VARCHAR(1000) DEFAULT '',

    -- Контакты (JSON)
    phones TEXT DEFAULT '[]',

    -- Классификация
    region_id INTEGER DEFAULT 0,
    city_id INTEGER DEFAULT 0,
    type_ad INTEGER DEFAULT 0,
    section_id INTEGER DEFAULT 0,
    category_id INTEGER DEFAULT 0,
    agent INTEGER DEFAULT 0,

    -- Временные метки
    first_seen_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_sent_at TIMESTAMP,
    api_updated_at VARCHAR(100) DEFAULT '',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_estates_id ON estates(id);
CREATE INDEX IF NOT EXISTS idx_estates_last_seen ON estates(last_seen_at);
CREATE INDEX IF NOT EXISTS idx_estates_region ON estates(region_id);
CREATE INDEX IF NOT EXISTS idx_estates_api_updated ON estates(api_updated_at);
CREATE INDEX IF NOT EXISTS idx_estates_first_seen ON estates(first_seen_at);
