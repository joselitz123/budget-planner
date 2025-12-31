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
    echo -e "${YELLOW}Warning: Go is not installed${NC}"
    echo "Backend may not start properly. Please install Go 1.23+"
fi

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

# 8. Start backend
echo ""
echo -e "${YELLOW}Starting backend server...${NC}"
cd backend

# Check if air is installed for hot reload
if command -v air &> /dev/null; then
    air > ../logs/backend.log 2>&1 &
    BACKEND_PID=$!
    echo -e "${GREEN}✓ Backend started with air (hot reload)${NC}"
else
    go run ./cmd/api > ../logs/backend.log 2>&1 &
    BACKEND_PID=$!
    echo -e "${GREEN}✓ Backend started with go run${NC}"
fi

cd ..

# Save backend PID
echo $BACKEND_PID > logs/backend.pid

# 9. Start frontend
echo ""
echo -e "${YELLOW}Starting frontend server...${NC}"
cd frontend
npm run dev > ../logs/frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

# Save frontend PID
echo $FRONTEND_PID > logs/frontend.pid

echo -e "${GREEN}✓ Frontend started${NC}"

# 10. Wait for services to be ready
echo ""
echo -e "${YELLOW}Waiting for services to start...${NC}"
sleep 5

# 11. Display status
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Budget Planner is running!${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "  Frontend: ${GREEN}http://localhost:5173${NC}"
echo -e "  Backend:  ${GREEN}http://localhost:8080${NC}"
echo ""
echo -e "Logs:"
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
