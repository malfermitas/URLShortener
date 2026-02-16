CREATE TABLE IF NOT EXISTS url_hit_events (
    id BIGSERIAL PRIMARY KEY,
    url_id VARCHAR(50) NOT NULL,
    user_agent VARCHAR(500),
    ip VARCHAR(45),
    country_code VARCHAR(10),
    referrer VARCHAR(500),
    device_type VARCHAR(50),
    os VARCHAR(50),
    browser VARCHAR(50),
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_url_hit_events_url_id ON url_hit_events(url_id);
CREATE INDEX IF NOT EXISTS idx_url_hit_events_timestamp ON url_hit_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_url_hit_events_country_code ON url_hit_events(country_code);
