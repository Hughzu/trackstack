UPDATE refills
SET season = NULL
WHERE season IS NOT NULL;
