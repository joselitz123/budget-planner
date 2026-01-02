#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# PIDs for cleanup
BACKEND_PID=0
FRONTEND_PID=0

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Budget Planner - Startup Script${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to cleanup on exit
cleanup() {
    echo ""
    echo -e "${YELLOW}Stopping services...${NC}"

    if [ $BACKEND_PID -ne 0 ]; then
        kill $BACKEND_PID 2>/dev/null || true
        echo -e "${GREEN}✓ Backend stopped${NC}"
    fi

    if [ $FRONTEND_PID -ne 0 ]; then
        kill $FRONTEND_PID 2>/dev/null || true
        echo -e "${GREEN}✓ Frontend stopped${NC}"
    fi

    # Stop Docker containers
    echo -e "${YELLOW}Stopping PostgreSQL...${NC}"
    docker-compose down 2>/dev/null || true
    echo -e "${GREEN}✓ PostgreSQL stopped${NC}"

    # Clean up PID files
    rm -f logs/backend.pid logs/frontend.pid

    echo -e "${GREEN}✓ All services stopped${NC}"
    exit 0
}

# Trap signals for cleanup
trap cleanup EXIT INT TERM

# 1. Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed${NC}"
    echo "Please install Docker first: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo -e "${RED}Error: Docker is not running${NC}"
    echo "Please start Docker and try again."
    exit 1
fi
echo -e "${GREEN}✓ Docker is running${NC}"

# 2. Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo -e "${RED}Error: Node.js is not installed${NC}"
    echo "Please install Node.js 18+: https://nodejs.org/"
    exit 1
fi
echo -e "${GREEN}✓ Node.js is installed$(node --version)${NC}"

# 3. Start PostgreSQL
echo ""
echo -e "${YELLOW}Starting PostgreSQL...${NC}"
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
for i in {1..30}; do
    if docker exec budget-planner-postgres pg_isready -U budgetuser -d budgetdb &> /dev/null; then
        echo -e "${GREEN}✓ PostgreSQL is ready${NC}"
        break
    fi
    sleep 1
done

# 4. Create .env files if they don't exist
echo ""
echo -e "${YELLOW}Setting up environment files...${NC}"

if [ ! -f backend/.env ]; then
    cp backend/.env.example backend/.env
    echo -e "${GREEN}✓ Created backend/.env${NC}"
else
    echo -e "${GREEN}✓ backend/.env already exists${NC}"
fi

if [ ! -f frontend/.env ]; then
    cp frontend/.env.example frontend/.env
    echo -e "${GREEN}✓ Created frontend/.env${NC}"
else
    echo -e "${GREEN}✓ frontend/.env already exists${NC}"
fi

# 5. Check Go for backend
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Backend may not start properly. Please install Go 1.23+"
    exit 1
fi
echo -e "${GREEN}✓ Go is installed$(go version)${NC}"

# Set up custom GOPATH to avoid permission issues
export GOPATH="$PWD/.go"
export GOMODCACHE="$GOPATH/pkg/mod"
mkdir -p "$GOMODCACHE"
echo -e "${GREEN}✓ GOPATH set to $GOPATH${NC}"

# 6. Install frontend dependencies if needed
echo ""
if [ ! -d frontend/node_modules ]; then
    echo -e "${YELLOW}Installing frontend dependencies...${NC}"
    cd frontend
    npm install --silent
    cd ..
    echo -e "${GREEN}✓ Dependencies installed${NC}"
else
    echo -e "${GREEN}✓ Frontend dependencies already installed${NC}"
fi

# 7. Create logs directory
mkdir -p logs

# 8. Clear caches to ensure environment variables are picked up
echo ""
echo -e "${YELLOW}Clearing caches...${NC}"
cd frontend

# Clear Vite cache
if [ -d .vite ]; then
    rm -rf .vite
    echo -e "${GREEN}✓ Cleared Vite cache${NC}"
fi

# Clear SvelteKit build cache (but keep tsconfig.json)
if [ -d .svelte-kit ]; then
    # Only clear specific cache directories, not the entire .svelte-kit
    rm -rf .svelte-kit/output .svelte-kit/server 2>/dev/null || true
    echo -e "${GREEN}✓ Cleared SvelteKit cache${NC}"
fi

# Clear node_modules/.vite cache
if [ -d node_modules/.vite ]; then
    rm -rf node_modules/.vite
    echo -e "${GREEN}✓ Cleared node_modules/.vite cache${NC}"
fi

cd ..

# 9. Kill existing backend process on port 8080 if running
echo ""
echo -e "${YELLOW}Checking for existing backend process on port 8080...${NC}"
if lsof -ti:8080 &> /dev/null; then
    echo -e "${YELLOW}Killing existing backend process on port 8080...${NC}"
    lsof -ti:8080 | xargs kill -9 2>/dev/null || true
    sleep 1
    echo -e "${GREEN}✓ Killed existing backend process${NC}"
else
    echo -e "${GREEN}✓ Port 8080 is available${NC}"
fi

# 10. Start backend
echo ""
echo -e "${YELLOW}Starting backend server...${NC}"
cd backend

# Check if air is installed for hot reload
if command -v air &> /dev/null; then
    echo -e "${BLUE}Using air for hot reload...${NC}"
    air 2>&1 | tee ../logs/backend.log &
    BACKEND_PID=$!
    echo -e "${GREEN}✓ Backend started with air (PID: $BACKEND_PID)${NC}"
else
    echo -e "${BLUE}Using go run...${NC}"
    go run ./cmd/api 2>&1 | tee ../logs/backend.log &
    BACKEND_PID=$!
    echo -e "${GREEN}✓ Backend started with go run (PID: $BACKEND_PID)${NC}"
fi

cd ..

# Save backend PID
echo $BACKEND_PID > logs/backend.pid

# 10. Start frontend
echo ""
echo -e "${YELLOW}Starting frontend server...${NC}"
cd frontend
npm run dev 2>&1 | tee ../logs/frontend.log &
FRONTEND_PID=$!
cd ..

# Save frontend PID
echo $FRONTEND_PID > logs/frontend.pid

echo -e "${GREEN}✓ Frontend started (PID: $FRONTEND_PID)${NC}"

# 11. Wait for services to be ready
echo ""
echo -e "${YELLOW}Waiting for services to start...${NC}"
sleep 5

# 12. Health checks
echo ""
echo -e "${YELLOW}Checking service health...${NC}"

# Check if backend process is still running
if ps -p $BACKEND_PID > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Backend process is running (PID: $BACKEND_PID)${NC}"
else
    echo -e "${RED}✗ Backend process is NOT running!${NC}"
    echo -e "${RED}Check logs/backend.log for errors${NC}"
fi

# Check if backend port is listening
if netstat -tlnp 2>/dev/null | grep -q ":8080 " || lsof -i :8080 2>/dev/null | grep -q LISTEN; then
    echo -e "${GREEN}✓ Backend is listening on port 8080${NC}"
else
    echo -e "${RED}✗ Backend is NOT listening on port 8080!${NC}"
    echo -e "${RED}The backend may have failed to start${NC}"
fi

# Check if frontend process is still running
if ps -p $FRONTEND_PID > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Frontend process is running (PID: $FRONTEND_PID)${NC}"
else
    echo -e "${RED}✗ Frontend process is NOT running!${NC}"
    echo -e "${RED}Check logs/frontend.log for errors${NC}"
fi

# Check if frontend port is listening
if netstat -tlnp 2>/dev/null | grep -q ":5173 " || lsof -i :5173 2>/dev/null | grep -q LISTEN; then
    echo -e "${GREEN}✓ Frontend is listening on port 5173${NC}"
else
    echo -e "${YELLOW}⚠ Frontend may still be starting up...${NC}"
fi

# 13. Display status
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Budget Planner is running!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "  Frontend: ${GREEN}http://localhost:5173${NC}"
echo -e "  Backend:  ${GREEN}http://localhost:8080${NC}"
echo ""
echo -e "Process IDs:"
echo -e "  Backend:  ${YELLOW}$BACKEND_PID${NC}"
echo -e "  Frontend: ${YELLOW}$FRONTEND_PID${NC}"
echo ""
echo -e "Logs (output also shown above):"
echo -e "  Backend:  ${YELLOW}logs/backend.log${NC}"
echo -e "  Frontend: ${YELLOW}logs/frontend.log${NC}"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo ""

# 12. Try to open browser (if on macOS or Linux with xdg-open)
if [[ "$OSTYPE" == "darwin"* ]]; then
    open http://localhost:5173 2>/dev/null || true
elif command -v xdg-open &> /dev/null; then
    xdg-open http://localhost:5173 2>/dev/null || true
fi

# Keep script running until interrupted
wait
