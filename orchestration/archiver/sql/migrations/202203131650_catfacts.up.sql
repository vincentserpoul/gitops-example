CREATE TABLE IF NOT EXISTS happycat_facts (
    id              uuid PRIMARY KEY,
    fact            TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);