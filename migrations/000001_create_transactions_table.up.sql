
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    amount NUMERIC NOT NULL,
    category TEXT,
    created_at TIMESTAMP DEFAULT NOW()
)