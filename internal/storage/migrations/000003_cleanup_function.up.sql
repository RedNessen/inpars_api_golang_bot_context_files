-- Функция для автоочистки старых объявлений
CREATE OR REPLACE FUNCTION cleanup_old_estates(days_threshold INTEGER)
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    -- Удалить объявления, которые не видели N+ дней
    DELETE FROM estates
    WHERE last_seen_at < NOW() - (days_threshold || ' days')::INTERVAL;

    GET DIAGNOSTICS deleted_count = ROW_COUNT;
    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
