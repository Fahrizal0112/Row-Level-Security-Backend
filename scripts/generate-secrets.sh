#!/bin/bash
echo "Generating secure secrets for your application..."
echo ""

echo "JWT_SECRET (copy this to your .env file):"
echo "JWT_SECRET=$(openssl rand -base64 32)"
echo ""

echo "API_KEY (copy this to your .env file):"
echo "API_KEY=$(openssl rand -hex 16)"
echo ""

echo "Database password (use this for your DATABASE_URL):"
echo "DB_PASSWORD=$(openssl rand -base64 16)"
echo ""

echo "Example DATABASE_URL:"
echo "DATABASE_URL=postgres://username:\$(DB_PASSWORD)@localhost:5432/row_level_security_db?sslmode=disable"