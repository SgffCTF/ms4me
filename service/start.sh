#!/bin/sh

echo "Starting ms4me service..."
JWT_SECRET=$(openssl rand -hex 32) docker compose up --build -d