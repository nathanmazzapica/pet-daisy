CREATE TABLE users (
                        user_id TEXT PRIMARY KEY,
                        pets INTEGER NOT NULL DEFAULT 0,
                        display_name TEXT NOT NULL,
                        sync_code TEXT NOT NULL DEFAULT 'NEEDCODEPLS',
                        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
