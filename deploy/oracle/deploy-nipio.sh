#!/usr/bin/env bash

set -euo pipefail

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required. Run deploy/oracle/oracle-first-boot.sh first."
  exit 1
fi

PUBLIC_IP="${1:-}"
if [ -z "$PUBLIC_IP" ]; then
  PUBLIC_IP="$(curl -fsSL https://api.ipify.org || true)"
fi

if [ -z "$PUBLIC_IP" ]; then
  echo "Could not detect public IP automatically."
  echo "Usage: ./deploy/oracle/deploy-nipio.sh <vm-public-ip>"
  exit 1
fi

DOMAIN="${DOMAIN:-${PUBLIC_IP}.nip.io}"
LETSENCRYPT_EMAIL="${LETSENCRYPT_EMAIL:-ops@example.com}"

if [ ! -f .env ]; then
  cp .env.example .env
fi

upsert_env() {
  local key="$1"
  local value="$2"

  if grep -q "^${key}=" .env; then
    sed -i "s|^${key}=.*|${key}=${value}|" .env
  else
    printf "%s=%s\n" "$key" "$value" >> .env
  fi
}

upsert_env "DOMAIN" "$DOMAIN"
upsert_env "BASE_URL" "https://${DOMAIN}"
upsert_env "LETSENCRYPT_EMAIL" "$LETSENCRYPT_EMAIL"
upsert_env "APP_BIND_ADDRESS" "127.0.0.1"
upsert_env "USE_MINIO" "false"

echo "Deploying FotoBoo with HTTPS domain: ${DOMAIN}"
docker compose --profile prod up -d --build
docker compose ps

echo "Done."
echo "Health check: curl -f https://${DOMAIN}/health"
