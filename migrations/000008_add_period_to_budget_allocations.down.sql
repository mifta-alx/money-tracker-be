BEGIN;

ALTER TABLE budget_allocations DROP CONSTRAINT IF EXISTS check_period_type;

ALTER TABLE budget_allocations DROP COLUMN IF EXISTS period;

COMMIT;