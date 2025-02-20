CREATE SCHEMA IF NOT EXISTS chatserver;

CREATE TABLE IF NOT EXISTS chatserver.channels (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    is_private BOOLEAN DEFAULT FALSE,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

ALTER TABLE chatserver.channels 
ALTER COLUMN updated_at SET DEFAULT now();
