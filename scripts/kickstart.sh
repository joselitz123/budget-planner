#!/bin/bash
# KICK-START helper script for Budget Planner project
# Usage: ./scripts/kickstart.sh

echo "🚀 KICK-START: Budget Planner Development Session"
echo "=================================================="
echo ""

# Check if beans is available
if ! command -v beans &> /dev/null; then
    echo "⚠️  Beans CLI not found in PATH"
    echo "   Try: export PATH=\$PATH:/home/vscode/go/bin"
    echo ""
fi

# Step 1: Beans context
echo "📋 Open Beans (by priority):"
echo "----------------------------"
if command -v beans &> /dev/null; then
    beans list --status todo 2>/dev/null || echo "No beans found"
else
    echo "Beans CLI not available"
fi
echo ""

# Step 2: Git status
echo "📂 Git Status:"
echo "-------------"
git status --short
echo ""

# Step 3: Recent commits
echo "📝 Recent Commits:"
echo "-----------------"
git log --oneline -5
echo ""

# Step 4: Environment check
echo "🔧 Environment Check:"
echo "--------------------"

# Check PostgreSQL
if docker ps &> /dev/null 2>&1; then
    if docker ps | grep -q postgres; then
        echo "✅ PostgreSQL running"
    else
        echo "⚠️  PostgreSQL not running"
        echo "   Start with: docker-compose up -d postgres"
    fi
else
    echo "⚠️  Docker not available"
fi

# Check current directory
echo ""
echo "📁 Current Directory: $(pwd)"
echo ""
echo "✅ Ready to code!"
echo ""
echo "Quick commands:"
echo "  - Frontend: cd frontend && npm run dev"
echo "  - Backend:  cd backend && air"
echo "  - Beans:    beans list"
