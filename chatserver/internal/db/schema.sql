CREATE SCHEMA IF NOT EXISTS chatserver;

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
