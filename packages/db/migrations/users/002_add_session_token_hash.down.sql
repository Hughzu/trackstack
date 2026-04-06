-- Rollback is currently disabled in CI because Atlas Community does not support migrate down.
-- This file stays intentionally non-destructive so `atlas migrate apply` does not sabotage itself.
SELECT 1;
