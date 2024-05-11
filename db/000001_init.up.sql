CREATE TABLE commands (
    id SERIAL PRIMARY KEY,
    script TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    executed_at TIMESTAMP,
    stdout TEXT,
    stderr TEXT
)