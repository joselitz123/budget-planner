#!/bin/bash

set -e

echo "🚀 Setting up Budget Planner development environment..."

# Install Go dependencies
echo "📦 Installing Go dependencies..."
cd /workspace/backend
if [ -f "go.mod" ]; then
    go mod download
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    go install github.com/air-verse/air@latest
else
    echo "⚠️  go.mod not found, skipping Go dependencies"
fi

# Install Node.js dependencies for frontend
echo "📦 Installing Node.js dependencies..."
cd /workspace/frontend
if [ -f "package.json" ]; then
    npm install
    npm install -D @sveltejs/adapter-node
else
    echo "⚠️  package.json not found, skipping Node.js dependencies"
fi

# Setup git hooks if .git exists
if [ -d "/workspace/.git" ]; then
    echo "🔧 Setting up git hooks..."
    git config core.hooksPath .githooks
    chmod +x .githooks/*
fi

# Create environment files if they don't exist
echo "📝 Creating environment files..."
cd /workspace

# Backend .env
if [ ! -f "backend/.env" ]; then
    cat > backend/.env << EOF
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=budgetuser
DB_PASSWORD=budgetpass
DB_NAME=budgetdb
DB_SSL_MODE=disable

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# Server
PORT=8080
ENVIRONMENT=development
ALLOWED_ORIGINS=http://localhost:5173

# Clerk Authentication
CLERK_SECRET_KEY=
CLERK_PUBLISHABLE_KEY=
CLERK_JWT_KEY=

# Sync
SYNC_BATCH_SIZE=50
SYNC_RETRY_ATTEMPTS=3
SYNC_RETRY_DELAY=5s

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
EOF
    echo "✅ Created backend/.env"
fi

# Frontend .env
if [ ! -f "frontend/.env" ]; then
    cat > frontend/.env << EOF
# Public
PUBLIC_CLERK_PUBLISHABLE_KEY=
PUBLIC_CLERK_SIGN_IN_URL=/sign-in
PUBLIC_CLERK_SIGN_UP_URL=/sign-up
PUBLIC_CLERK_AFTER_SIGN_IN_URL=/dashboard
PUBLIC_CLERK_AFTER_SIGN_UP_URL=/dashboard

# API
PUBLIC_API_URL=http://localhost:8080/api

# PWA
PUBLIC_APP_NAME="Budget Planner"
PUBLIC_APP_SHORT_NAME="Budget"
PUBLIC_APP_DESCRIPTION="Offline-first budget planning application"
PUBLIC_APP_THEME_COLOR="#3b82f6"
PUBLIC_APP_BACKGROUND_COLOR="#ffffff"

# Sync
PUBLIC_SYNC_INTERVAL=30000 # 30 seconds
PUBLIC_OFFLINE_RETRY_DELAY=5000 # 5 seconds
PUBLIC_MAX_OFFLINE_OPERATIONS=100
EOF
    echo "✅ Created frontend/.env"
fi

# Make scripts executable
chmod +x .devcontainer/post-create.sh

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
until pg_isready -h postgres -U budgetuser -d budgetdb; do
    echo "Waiting for PostgreSQL..."
    sleep 2
done

# Wait for Redis to be ready
echo "⏳ Waiting for Redis to be ready..."
until redis-cli -h redis ping | grep -q "PONG"; do
    echo "Waiting for Redis..."
    sleep 2
done

echo "✅ PostgreSQL and Redis are ready!"

# Run database migrations if sqlc is configured
cd /workspace/backend
if [ -f "sqlc.yaml" ] && [ -f "sql/schema" ]; then
    echo "🗄️  Setting up database schema..."
    # This would run sqlc generate and apply migrations
    # For now, just create a placeholder
    echo "📋 Database schema setup placeholder"
fi

# Create development certificates for HTTPS if needed
echo "🔐 Setting up development certificates..."
mkdir -p /workspace/certs
if [ ! -f "/workspace/certs/localhost.key" ]; then
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout /workspace/certs/localhost.key \
        -out /workspace/certs/localhost.crt \
        -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost" 2>/dev/null
    echo "✅ Created development certificates"
fi

echo "🎉 Development environment setup complete!"
echo ""
echo "Next steps:"
echo "1. Update the .env files with your Clerk API keys"
echo "2. Start the development servers:"
echo "   - Backend: cd backend && air"
echo "   - Frontend: cd frontend && npm run dev"
echo "3. Access the application at http://localhost:5173"
echo ""
echo "Happy coding! 💻"
