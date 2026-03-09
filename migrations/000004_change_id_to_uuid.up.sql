CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE users ADD COLUMN new_id UUID DEFAULT uuid_generate_v4();
ALTER TABLE auth_providers ADD column new_user_id UUID;

UPDATE auth_providers SET new_user_id = (SELECT new_id FROM users WHERE users.id = auth_providers.user_id);

ALTER TABLE auth_providers DROP CONSTRAINT IF EXISTS auth_providers_user_id_fkey;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey CASCADE;

ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN new_id TO id;
ALTER TABLE users ADD PRIMARY KEY (id);

ALTER TABLE auth_providers DROP COLUMN user_id;
ALTER TABLE auth_providers RENAME COLUMN new_user_id TO user_id;
ALTER TABLE auth_providers ADD CONSTRAINT auth_providers_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 