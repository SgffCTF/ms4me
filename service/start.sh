#!/bin/sh

set -e

echo "Generating JWT secret..."
JWT_SECRET=$(openssl rand -hex 32)

echo "Writing .env..."
cat <<EOF > .env
JWT_SECRET=$JWT_SECRET
EOF

echo "Starting ms4me services..."
docker compose --env-file .env up --build -d
