BEGIN;

ALTER TABLE accounts ADD COLUMN icon VARCHAR(50);
UPDATE accounts SET icon = 'wallet' WHERE icon IS NULL;
ALTER TABLE accounts ALTER COLUMN icon SET NOT NULL;

DROP TABLE IF EXISTS transactions CASCADE;
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    title VARCHAR(255),
    type VARCHAR(20), -- 'expense' atau 'income'
    amount BIGINT NOT NULL DEFAULT 0,
    date DATE DEFAULT CURRENT_DATE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
COMMIT;