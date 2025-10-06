INSERT INTO tenants (name, domain) VALUES 
('Default Tenant', 'default.local'),
('Company A', 'company-a.com'),
('Company B', 'company-b.com');

INSERT INTO users (email, password, name, tenant_id, role) VALUES 
('admin@default.local', '$2a$14$example_hashed_password', 'Admin User', 1, 'admin'),
('user1@company-a.com', '$2a$14$example_hashed_password', 'User One', 2, 'user'),
('user2@company-b.com', '$2a$14$example_hashed_password', 'User Two', 3, 'user');