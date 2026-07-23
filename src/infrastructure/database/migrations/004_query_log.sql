-- 004_query_log.sql — Audit log for NL→SQL generated queries (NFR-1)

CREATE TABLE IF NOT EXISTS query_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    question TEXT NOT NULL,
    generated_sql TEXT,
    success INTEGER NOT NULL DEFAULT 0,
    error_message TEXT,
    execution_time_ms INTEGER,
    created_at TEXT NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE INDEX IF NOT EXISTS idx_query_log_created ON query_log(created_at);
