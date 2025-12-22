-- Database initialization script for Budget Planner
-- Run automatically when PostgreSQL container starts

-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_user_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    currency VARCHAR(3) DEFAULT 'PHP',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted BOOLEAN DEFAULT FALSE
);

-- Categories Table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(10),
    color VARCHAR(7) DEFAULT '#3498db',
    is_system BOOLEAN DEFAULT FALSE,
    default_limit DECIMAL(12, 2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted BOOLEAN DEFAULT FALSE
);

-- Budgets Table
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100),
    month DATE NOT NULL,
    total_limit DECIMAL(12, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted BOOLEAN DEFAULT FALSE,
    UNIQUE(user_id, month)
);

-- Budget Categories Table
CREATE TABLE budget_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID REFERENCES budgets(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    limit_amount DECIMAL(12, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(budget_id, category_id)
);

-- Payment Methods Table
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    last_four VARCHAR(4),
    brand VARCHAR(50),
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    credit_limit DECIMAL(12, 2),
    current_balance DECIMAL(12, 2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted BOOLEAN DEFAULT FALSE
);

-- Transactions Table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    budget_id UUID REFERENCES budgets(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    payment_method_id UUID REFERENCES payment_methods(id) ON DELETE SET NULL,
    amount DECIMAL(12, 2) NOT NULL,
    type VARCHAR(10) DEFAULT 'expense',
    is_transfer BOOLEAN DEFAULT FALSE,
    transfer_to_account_id UUID,
    description VARCHAR(255),
    transaction_date DATE NOT NULL DEFAULT CURRENT_DATE,
    is_recurring BOOLEAN DEFAULT FALSE,
    recurrence_pattern JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted BOOLEAN DEFAULT FALSE
);

-- Reflections Table
CREATE TABLE reflections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    budget_id UUID REFERENCES budgets(id) ON DELETE CASCADE,
    overall_rating INTEGER CHECK (overall_rating >= 1 AND overall_rating <= 10),
    is_private BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted BOOLEAN DEFAULT FALSE,
    UNIQUE(user_id, budget_id)
);

-- Reflection Questions Table
CREATE TABLE reflection_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reflection_id UUID REFERENCES reflections(id) ON DELETE CASCADE,
    sequence INTEGER DEFAULT 0,
    question_id VARCHAR(50) NOT NULL,
    question_text TEXT NOT NULL,
    answer TEXT,
    mood VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Reflection Templates Table
CREATE TABLE reflection_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Template Questions Table
CREATE TABLE template_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID REFERENCES reflection_templates(id) ON DELETE CASCADE,
    question_id VARCHAR(50) NOT NULL,
    question_text TEXT NOT NULL,
    type VARCHAR(20) DEFAULT 'text',
    is_required BOOLEAN DEFAULT TRUE,
    has_mood_tracking BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(template_id, question_id)
);

-- Share Invitations Table
CREATE TABLE share_invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID REFERENCES budgets(id) ON DELETE CASCADE,
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    recipient_email VARCHAR(255) NOT NULL,
    permission VARCHAR(10) NOT NULL CHECK (permission IN ('view', 'edit')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined', 'expired')),
    expires_at TIMESTAMPTZ DEFAULT (NOW() + INTERVAL '7 days'),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Share Access Table
CREATE TABLE share_access (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID REFERENCES budgets(id) ON DELETE CASCADE,
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    shared_with_id UUID REFERENCES users(id) ON DELETE CASCADE,
    permission VARCHAR(10) NOT NULL CHECK (permission IN ('view', 'edit')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(budget_id, shared_with_id)
);

-- Sync Table
CREATE TABLE sync_operations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    table_name VARCHAR(50) NOT NULL,
    record_id UUID NOT NULL,
    operation VARCHAR(10) NOT NULL,
    local_data JSONB,
    server_data JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    error_message TEXT,
    attempt_count INTEGER DEFAULT 0,
    last_attempt_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Activity Log Table
CREATE TABLE activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50),
    resource_id UUID,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Insert default system categories
INSERT INTO categories (id, user_id, name, icon, color, is_system, default_limit, created_at, updated_at) VALUES
(uuid_generate_v4(), NULL, 'Food', '🥗', '#e74c3c', TRUE, 10000.00, NOW(), NOW()),
(uuid_generate_v4(), NULL, 'Transport', '🚗', '#3498db', TRUE, 5000.00, NOW(), NOW()),
(uuid_generate_v4(), NULL, 'Housing', '🏠', '#2ecc71', TRUE, 20000.00, NOW(), NOW()),
(uuid_generate_v4(), NULL, 'Entertainment', '🎬', '#9b59b6', TRUE, 3000.00, NOW(), NOW()),
(uuid_generate_v4(), NULL, 'Healthcare', '🏥', '#e67e22', TRUE, 5000.00, NOW(), NOW()),
(uuid_generate_v4(), NULL, 'Education', '📚', '#1abc9c', TRUE, 5000.00, NOW(), NOW()),
(uuid_generate_v4(), NULL, 'Shopping', '🛍️', '#f39c12', TRUE, 8000.00, NOW(), NOW()),
(uuid_generate_v4(), NULL, 'Utilities', '💡', '#34495e', TRUE, 7000.00, NOW(), NOW()),
(uuid_generate_v4(), NULL, 'Others', '📦', '#95a5a6', TRUE, 4000.00, NOW(), NOW());

-- Insert default reflection template
INSERT INTO reflection_templates (id, name, is_active, version, created_at, updated_at) VALUES
(uuid_generate_v4(), 'Monthly Budget Reflection', TRUE, 1, NOW(), NOW());

-- Insert template questions
INSERT INTO template_questions (id, template_id, question_id, question_text, type, is_required, has_mood_tracking, sort_order, created_at, updated_at)
SELECT 
    uuid_generate_v4(),
    (SELECT id FROM reflection_templates WHERE name = 'Monthly Budget Reflection'),
    q.id,
    q.question,
    'text',
    TRUE,
    TRUE,
    q.sort_order,
    NOW(),
    NOW()
FROM (VALUES
    ('Q1', 'How satisfied are you with your spending this month?', 1),
    ('Q2', 'Did you stay within your budget for most categories?', 2),
    ('Q3', 'What was your biggest financial win this month?', 3),
    ('Q4', 'What spending surprised you the most?', 4),
    ('Q5', 'What would you do differently next month?', 5),
    ('Q6', 'How do you feel about your financial progress overall?', 6)
) AS q(id, question, sort_order);

-- Create Indexes
CREATE INDEX idx_transactions_user_month ON transactions(user_id, DATE_TRUNC('month', transaction_date));
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_budgets_month ON budgets(month) WHERE NOT deleted;
CREATE INDEX idx_categories_user_active ON categories(user_id) WHERE NOT deleted AND NOT is_system;
CREATE INDEX idx_share_access_shared_with ON share_access(shared_with_id) WHERE permission = 'edit';
CREATE INDEX idx_share_invitations_email ON share_invitations(recipient_email);
CREATE INDEX idx_sync_operations_created ON sync_operations(created_at DESC);
CREATE INDEX idx_sync_operations_attempt ON sync_operations(attempt_count) WHERE status = 'failed';

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
DO $$
DECLARE
    t text;
BEGIN
    FOR t IN 
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name IN (
            'users', 'categories', 'budgets', 'budget_categories', 
            'payment_methods', 'transactions', 'reflections', 
            'reflection_questions', 'reflection_templates', 'template_questions',
            'share_invitations', 'sync_operations'
        )
    LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS update_%s_updated_at ON %I', t, t);
        EXECUTE format('CREATE TRIGGER update_%s_updated_at BEFORE UPDATE ON %I FOR EACH ROW EXECUTE FUNCTION update_updated_at_column()', t, t);
    END LOOP;
END;
$$;
