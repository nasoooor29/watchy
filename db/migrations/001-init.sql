CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS cookies (
    uuid TEXT PRIMARY KEY,
    userid TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- Your own internal application ID
    title TEXT NOT NULL,
    mal_id INTEGER, -- For jikan/mal
);

CREATE INDEX IF NOT EXISTS idx_cookies_userid ON cookies(userid);
