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
    authored_at TIMESTAMP NOT NULL
);

ALTER TABLE chatserver.channels 
ALTER COLUMN updated_at SET DEFAULT now();
