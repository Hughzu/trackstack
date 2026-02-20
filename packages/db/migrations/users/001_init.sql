CREATE TABLE users (
  id TEXT PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  session_version INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  last_login_at TEXT
);

CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
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

CREATE INDEX sessions_user_id_idx ON sessions(user_id);
CREATE INDEX sessions_expires_at_idx ON sessions(expires_at);
CREATE INDEX sessions_absolute_expires_at_idx ON sessions(absolute_expires_at);
