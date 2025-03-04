CREATE SCHEMA IF NOT EXISTS chatserver;

CREATE TABLE IF NOT EXISTS chatserver.channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(24) NOT NULL UNIQUE,
    description TEXT,
    is_private BOOLEAN DEFAULT FALSE,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS chatserver.chat_messages (
    id SERIAL PRIMARY KEY,
    owner_id UUID NOT NULL,
    channel VARCHAR(24) NOT NULL,  -- Stores the channel name directly
    message TEXT NOT NULL,
    authored_at TIMESTAMP NOT NULL,
    search_vector TSVECTOR GENERATED ALWAYS AS (to_tsvector('english', message)) STORED
);

CREATE INDEX IF NOT EXISTS chat_messages_search_idx ON chatserver.chat_messages USING GIN(search_vector);

CREATE TABLE IF NOT EXISTS chatserver.chat_sessions (
    id SERIAL PRIMARY KEY,
    owner_id UUID NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    duration INTERVAL GENERATED ALWAYS AS (end_time - start_time) STORED
);
