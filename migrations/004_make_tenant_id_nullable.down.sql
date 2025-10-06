-- Revert tenant_id to NOT NULL (this will fail if there are NULL values)
ALTER TABLE users ALTER COLUMN tenant_id SET NOT NULL;

-- Revert RLS policy
DROP POLICY IF EXISTS user_tenant_isolation_policy ON users;

CREATE POLICY user_tenant_isolation_policy ON users
    FOR ALL
    TO PUBLIC
    USING (tenant_id = COALESCE(current_setting('app.current_tenant_id', true)::INTEGER, tenant_id));