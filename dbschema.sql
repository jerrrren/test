DROP TABLE IF EXISTS users;
CREATE TABLE users (
    uid		    SERIAL PRIMARY KEY,
    name		TEXT,
    password	TEXT,
    token	TEXT,
    refresh_token TEXT,
    user_type TEXT
);