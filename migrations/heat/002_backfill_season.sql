UPDATE refills
SET season = (
  CAST(strftime('%Y', date) AS INTEGER)
  - CASE WHEN CAST(strftime('%m', date) AS INTEGER) < 9 THEN 1 ELSE 0 END
) || '-' || (
  CAST(strftime('%Y', date) AS INTEGER)
  - CASE WHEN CAST(strftime('%m', date) AS INTEGER) < 9 THEN 1 ELSE 0 END
  + 1
)
WHERE season IS NULL;
