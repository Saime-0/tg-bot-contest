CREATE TABLE members
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chat_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status INTEGER NOT NULL DEFAULT 0,
    inviter_id INTEGER NOT NULL DEFAULT 0,
    ignore_in_ticket_counting INTEGER NOT NULL,
    in_ticket_id INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    is_bot INTEGER NOT NULL DEFAULT FALSE,
    first_name TEXT NOT NULL,
    username TEXT NOT NULL DEFAULT '',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chats(
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    username TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tickets(
    number INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    contest_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE contests(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    creator_id INTEGER,
    chat_id INTEGER NOT NULL,
    topic_id INTEGER NOT NULL DEFAULT 0,
    keyword TEXT NOT NULL,
    multiplicity INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP NULL
);

CREATE TABLE metadata (
    key TEXT PRIMARY KEY,
    value TEXT
);