BEGIN;

ALTER TABLE budget_allocations ADD COLUMN period VARCHAR(20) NOT NULL DEFAULT 'monthly';

ALTER TABLE budget_allocations ADD CONSTRAINT check_period_type CHECK (period IN ('weekly', 'monthly', 'yearly'));

COMMIT;