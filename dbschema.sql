DROP TABLE IF EXISTS users;
CREATE TABLE users (
    uid		    SERIAL PRIMARY KEY,
    name		TEXT,
    password	TEXT,
    token	TEXT,
    refresh_token TEXT,
    user_type TEXT
);

DROP TABLE IF EXISTS chats;
CREATE TABLE chats (
    messageID SERIAL PRIMARY KEY,
    user_id_1 INT, 
    user_id_2 INT,
    body TEXT,
    messageTime TEXT
);


DROP TABLE IF EXISTS posts;
CREATE TABLE posts(
    id SERIAL PRIMARY KEY,
    field TEXT NOT NULL,
    name TEXT NOT NULL,
    intro TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    participants TEXT[]
);

DROP TABLE IF EXISTS singleusers;
CREATE TABLE singleusers(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    commitment INTEGER,
    location TEXT,
    filledinfo BOOLEAN NOT NULL DEFAULT false
);

DROP TABLE IF EXISTS pairedusers;
CREATE TABLE pairedusers(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    partner TEXT NOT NULL UNIQUE
);