CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tenant_id INTEGER NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_tenant_id ON posts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts(deleted_at);

ALTER TABLE posts ENABLE ROW LEVEL SECURITY;

CREATE POLICY post_tenant_isolation_policy ON posts
    FOR SELECT
    TO PUBLIC
    USING (
        tenant_id = COALESCE(current_setting('app.current_tenant_id', true)::INTEGER, tenant_id)
        AND (
            is_public = true 
            OR user_id = COALESCE(current_setting('app.current_user_id', true)::INTEGER, user_id)
        )
    );

CREATE POLICY post_create_policy ON posts
    FOR INSERT
    TO PUBLIC
    WITH CHECK (
        tenant_id = COALESCE(current_setting('app.current_tenant_id', true)::INTEGER, tenant_id)
    );

CREATE POLICY post_update_policy ON posts
    FOR UPDATE
    TO PUBLIC
    USING (
        tenant_id = COALESCE(current_setting('app.current_tenant_id', true)::INTEGER, tenant_id)
        AND user_id = COALESCE(current_setting('app.current_user_id', true)::INTEGER, user_id)
    );

CREATE POLICY post_delete_policy ON posts
    FOR DELETE
    TO PUBLIC
    USING (
        tenant_id = COALESCE(current_setting('app.current_tenant_id', true)::INTEGER, tenant_id)
        AND user_id = COALESCE(current_setting('app.current_user_id', true)::INTEGER, user_id)
    );