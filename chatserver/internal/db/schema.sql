CREATE SCHEMA IF NOT EXISTS chatserver;

-- Update timestamp trigger function (must be declared before any trigger uses it)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Tables
CREATE TABLE IF NOT EXISTS chatserver.server_instances (
    id SERIAL PRIMARY KEY,
    server_id UUID NOT NULL UNIQUE,
    server_name VARCHAR(36) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS chatserver.rate_limiter (
    id SERIAL PRIMARY KEY,
    owner_id VARCHAR(36) NOT NULL UNIQUE,
    message_limit INT NOT NULL DEFAULT 10,
    window_seconds INT NOT NULL DEFAULT 60,
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS chatserver.channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(24) NOT NULL UNIQUE,
    description TEXT,
    is_private BOOLEAN DEFAULT FALSE,
    owner_id VARCHAR(36) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS chatserver.chat_messages (
    id SERIAL PRIMARY KEY,
    cache_id BIGINT UNIQUE NOT NULL,  -- Stores Valkey's INCR value
    owner_id VARCHAR(36) NOT NULL,
    channel VARCHAR(24) NOT NULL,  -- Stores the channel name directly (not channel ID)
    message TEXT NOT NULL,
    authored_at TIMESTAMP NOT NULL,
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', message)) STORED
);

CREATE INDEX IF NOT EXISTS chat_messages_search_idx ON chatserver.chat_messages USING GIN(search_vector);

CREATE TABLE IF NOT EXISTS chatserver.chat_sessions (
    id SERIAL PRIMARY KEY,
    owner_id VARCHAR(36) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration INTERVAL GENERATED ALWAYS AS (end_time - start_time) STORED
);

CREATE TABLE IF NOT EXISTS chatserver.bans (
    id SERIAL PRIMARY KEY,
    owner_id VARCHAR(36) NOT NULL,
    banished_id VARCHAR(36) NOT NULL,
    start_time TIMESTAMP DEFAULT now(),
    reason VARCHAR(256),
    end_time TIMESTAMP NULL,
    duration INTERVAL GENERATED ALWAYS AS 
        (CASE WHEN end_time IS NOT NULL THEN end_time - start_time ELSE NULL END) STORED,
    pardoned BOOLEAN DEFAULT FALSE
);

-- Triggers to auto-update 'updated_at' timestamp
DROP TRIGGER IF EXISTS update_server_instances_updated_at ON chatserver.server_instances;
CREATE TRIGGER update_server_instances_updated_at
BEFORE UPDATE ON chatserver.server_instances
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_channels_updated_at ON chatserver.channels;
CREATE TRIGGER update_channels_updated_at
BEFORE UPDATE ON chatserver.channels
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_rate_limiter_updated_at ON chatserver.rate_limiter;
CREATE TRIGGER update_rate_limiter_updated_at
BEFORE UPDATE ON chatserver.rate_limiter
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
