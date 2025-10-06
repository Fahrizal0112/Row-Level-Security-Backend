-- Make tenant_id nullable to allow users without tenants
ALTER TABLE users ALTER COLUMN tenant_id DROP NOT NULL;

-- Update RLS policy to handle NULL tenant_id
DROP POLICY IF EXISTS user_tenant_isolation_policy ON users;

CREATE POLICY user_tenant_isolation_policy ON users
    FOR ALL
    TO PUBLIC
    USING (
        tenant_id IS NULL OR 
        tenant_id = COALESCE(current_setting('app.current_tenant_id', true)::INTEGER, tenant_id)
    );