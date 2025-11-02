CREATE TABLE IF NOT EXISTS user (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE,
    password_hash TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS manga (
    id TEXT PRIMARY KEY,
    title TEXT,
    author TEXT,
    genres TEXT, -- JSON array as text
    status TEXT,
    total_chapters INTEGER,
    description TEXT
);

CREATE TABLE IF NOT EXISTS user_progress (
    user_id TEXT,
    manga_id TEXT,
    current_chapter INTEGER,
    status TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, manga_id)
);
