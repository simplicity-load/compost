CREATE TABLE IF NOT EXISTS "users" (
    "id"    INTEGER NOT NULL,
    "name"  TEXT NOT NULL,
    PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "tasks" (
    "id"        INTEGER NOT NULL UNIQUE,
    "user_id"   INTEGER NOT NULL,
    "title"     TEXT,
    "body"      TEXT,
    "status"    TEXT NOT NULL,
    PRIMARY KEY("id" AUTOINCREMENT),
    FOREIGN KEY("user_id") REFERENCES "users" ("id")
);
