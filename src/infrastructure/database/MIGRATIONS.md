# SQLite migration lifecycle

Fresh databases apply embedded migrations in filename order. Each migration's SQL
and ledger entry commit together, so a failed migration rolls back and can be
retried; completed migrations are skipped on later boots. Demo seed SQL runs once.

Existing databases without `schema_migrations` are rejected without altering their
schema or data. Automatic legacy adoption/import is not supported. Render Free
uses a fresh disposable database, not the local `data/pos.db`; never upload the
local database or delete it to resolve a deployment error.
