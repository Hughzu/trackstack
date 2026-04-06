DROP INDEX IF EXISTS sessions_user_id_idx;
DROP INDEX IF EXISTS sessions_expires_at_idx;
DROP INDEX IF EXISTS sessions_absolute_expires_at_idx;

ALTER TABLE sessions RENAME TO sessions_legacy_without_token_hash;

CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  token_hash TEXT NOT NULL UNIQUE,
  created_at TEXT NOT NULL,
  expires_at TEXT NOT NULL,
  rotated_at TEXT NOT NULL,
  last_seen_at TEXT NOT NULL,
  absolute_expires_at TEXT NOT NULL,
  parent_id TEXT,
  revoked_at TEXT,
  user_agent_hash TEXT,
  ip_prefix TEXT,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO sessions (
  id,
  user_id,
  token_hash,
  created_at,
  expires_at,
  rotated_at,
  last_seen_at,
  absolute_expires_at,
  parent_id,
  revoked_at,
  user_agent_hash,
  ip_prefix
)
SELECT
  id,
  user_id,
  lower(hex(randomblob(32))),
  created_at,
  expires_at,
  rotated_at,
  last_seen_at,
  absolute_expires_at,
  parent_id,
  revoked_at,
  user_agent_hash,
  ip_prefix
FROM sessions_legacy_without_token_hash;

DROP TABLE sessions_legacy_without_token_hash;

CREATE INDEX sessions_user_id_idx ON sessions(user_id);
CREATE INDEX sessions_token_hash_idx ON sessions(token_hash);
CREATE INDEX sessions_expires_at_idx ON sessions(expires_at);
CREATE INDEX sessions_absolute_expires_at_idx ON sessions(absolute_expires_at);
