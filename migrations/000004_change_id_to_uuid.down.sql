ALTER TABLE users ADD COLUMN old_id SERIAL;
ALTER TABLE auth_providers ADD COLUMN old_user_id INT;

ALTER TABLE auth_providers DROP CONSTRAINT IF EXISTS auth_providers_user_id_fkey;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey CASCADE;

ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN old_id TO id;
ALTER TABLE users ADD PRIMARY KEY (id);

ALTER TABLE auth_providers DROP COLUMN user_id;
ALTER TABLE auth_providers RENAME COLUMN old_user_id TO user_id;
ALTER TABLE auth_providers ADD CONSTRAINT auth_providers_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;