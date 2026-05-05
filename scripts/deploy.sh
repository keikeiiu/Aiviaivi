#!/bin/bash
# AiliVili deployment script for fly.io
# Prerequisites: brew install flyctl (macOS) or curl -L https://fly.io/install.sh | sh
set -euo pipefail

echo "=== AiliVili Deployment ==="

# Check fly CLI
if ! command -v fly &>/dev/null; then
    echo "Installing flyctl..."
    curl -L https://fly.io/install.sh | sh
    export PATH="$HOME/.fly/bin:$PATH"
fi

# Check if app exists
if ! fly apps list 2>/dev/null | grep -q "ailivili"; then
    echo "Creating app..."
    fly apps create ailivili
fi

# Create managed PostgreSQL if needed
if ! fly pg list 2>/dev/null | grep -q "ailivili-db"; then
    echo "Creating PostgreSQL (this takes ~2min)..."
    fly pg create --name ailivili-db --region nrt --vm-size shared-cpu-1x --initial-cluster-size 1
    echo "Attaching PostgreSQL to app..."
    fly pg attach ailivili-db --app ailivili
fi

# Create Redis if needed
if ! fly redis list 2>/dev/null | grep -q "ailivili-redis"; then
    echo "Creating Redis..."
    fly redis create --name ailivili-redis --region nrt
fi

# Set secrets (prompt for JWT secret)
echo ""
read -rp "JWT_SECRET (enter for random): " JWT_SECRET
JWT_SECRET="${JWT_SECRET:-$(openssl rand -hex 32)}"
fly secrets set JWT_SECRET="$JWT_SECRET" --app ailivili
fly secrets set JWT_EXPIRES_MINUTES=60 --app ailivili
fly secrets set CORS_ORIGIN="https://ailivili.fly.dev" --app ailivili
fly secrets set STORAGE=local --app ailivili

# Deploy
echo ""
echo "=== Deploying ==="
fly deploy --app ailivili

echo ""
echo "=== Done ==="
echo "API: https://ailivili.fly.dev/api/v1/health"
echo "Metrics: https://ailivili.fly.dev/metrics"
echo ""
echo "Frontend: update frontend/utils/constants.ts with:"
echo "  API_BASE_URL = 'https://ailivili.fly.dev/api/v1'"
echo "  WS_BASE_URL = 'wss://ailivili.fly.dev/api/v1'"
