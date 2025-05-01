CREATE SCHEMA IF NOT EXISTS chatserver;

-- Trigger function to automatically update the 'updated_at' column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ====================================
-- Server Identity & Configuration
-- ====================================

-- Stores registered instances of the server (useful for multi-server deployments)
CREATE TABLE IF NOT EXISTS chatserver.server_instances (
    id SERIAL PRIMARY KEY,
    server_id UUID NOT NULL UNIQUE,         -- UUID representing this server instance
    server_name VARCHAR(36) NOT NULL UNIQUE, -- Human-readable identifier
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- Stores per-user rate limiting rules (set by admins)
CREATE TABLE IF NOT EXISTS chatserver.rate_limiter (
    id SERIAL PRIMARY KEY,
    owner_id VARCHAR(36) NOT NULL UNIQUE,    -- User ID
    message_limit INT NOT NULL DEFAULT 10,   -- Max allowed messages per window
    window_seconds INT NOT NULL DEFAULT 60,  -- Sliding window in seconds
    updated_at TIMESTAMP DEFAULT now()
);

-- ====================================
-- Channel Metadata
-- ====================================

-- Stores chat channels, public or private
CREATE TABLE IF NOT EXISTS chatserver.channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(24) NOT NULL UNIQUE,       -- Channel slug (e.g., 'general', 'tech')
    description TEXT,
    is_private BOOLEAN DEFAULT FALSE,       -- If true, not shown to all users
    owner_id VARCHAR(36) NOT NULL,          -- Creator of the channel
    sort_order INT NOT NULL DEFAULT 0,      -- Used for custom sorting in UI
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- ====================================
-- Messages
-- ====================================

-- Stores public chat messages (channel-based)
CREATE TABLE IF NOT EXISTS chatserver.chat_messages (
    id SERIAL PRIMARY KEY,
    cache_id BIGINT UNIQUE NOT NULL,        -- Unique message ID from Valkey
    owner_id VARCHAR(36) NOT NULL,
    channel VARCHAR(24) NOT NULL,           -- Denormalized channel name (not a foreign key)
    message TEXT NOT NULL,
    authored_at TIMESTAMP NOT NULL,
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', message)) STORED
);

-- GIN index for fast full-text search
CREATE INDEX IF NOT EXISTS chat_messages_search_idx ON chatserver.chat_messages USING GIN(search_vector);

-- Stores private messages between two users
CREATE TABLE IF NOT EXISTS chatserver.private_messages (
    id SERIAL PRIMARY KEY,
    cache_id BIGINT UNIQUE NOT NULL,        -- Unique message ID from Valkey
    owner_id VARCHAR(36) NOT NULL,          -- Sender ID
    username VARCHAR(64) NOT NULL,          -- Sender username (denormalized)
    recipient_id VARCHAR(36) NOT NULL,      -- Receiver ID
    recipient VARCHAR(64) NOT NULL,         -- Receiver username (denormalized)
    message TEXT NOT NULL,
    authored_at TIMESTAMP NOT NULL,
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', message)) STORED
);

CREATE INDEX IF NOT EXISTS private_messages_search_idx ON chatserver.private_messages USING GIN(search_vector);

-- ====================================
-- Session and Moderation Records
-- ====================================

-- Stores session duration metrics for users
CREATE TABLE IF NOT EXISTS chatserver.chat_sessions (
    id SERIAL PRIMARY KEY,
    owner_id VARCHAR(36) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration INTERVAL GENERATED ALWAYS AS (end_time - start_time) STORED
);

-- Stores ban records for moderation
CREATE TABLE IF NOT EXISTS chatserver.bans (
    id SERIAL PRIMARY KEY,
    owner_id VARCHAR(36) NOT NULL,          -- Admin/mod ID
    banished_id VARCHAR(36) NOT NULL,       -- Banned user's ID
    start_time TIMESTAMP DEFAULT now(),
    reason VARCHAR(256),
    end_time TIMESTAMP NULL,                -- NULL means indefinite ban
    duration INTERVAL GENERATED ALWAYS AS 
        (CASE WHEN end_time IS NOT NULL THEN end_time - start_time ELSE NULL END) STORED,
    pardoned BOOLEAN DEFAULT FALSE          -- If true, ban is forgiven
);

-- ====================================
-- Auto-Update Triggers
-- ====================================

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
