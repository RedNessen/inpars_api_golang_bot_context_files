-- Таблица истории изменений
CREATE TABLE IF NOT EXISTS estate_changes (
    id SERIAL PRIMARY KEY,
    estate_id BIGINT NOT NULL REFERENCES estates(id) ON DELETE CASCADE,

    field_name VARCHAR(50) NOT NULL,
    old_value TEXT DEFAULT '',
    new_value TEXT DEFAULT '',

    changed_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_changes_estate_id ON estate_changes(estate_id);
CREATE INDEX IF NOT EXISTS idx_changes_changed_at ON estate_changes(changed_at);
