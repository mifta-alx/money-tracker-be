BEGIN;

ALTER TABLE users ADD COLUMN IF NOT EXISTS last_sign_in_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE auth_providers ADD COLUMN new_id UUID DEFAULT uuid_generate_v4();
ALTER TABLE auth_providers DROP COLUMN id;
ALTER TABLE auth_providers RENAME COLUMN new_id TO id;
ALTER TABLE auth_providers ADD PRIMARY KEY (id);

CREATE INDEX IF NOT EXISTS idx_users_last_sign_in ON users(last_sign_in_at);

COMMIT;