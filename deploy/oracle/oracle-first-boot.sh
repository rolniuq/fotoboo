#!/usr/bin/env bash

set -euo pipefail

if ! command -v apt-get >/dev/null 2>&1; then
  echo "This script currently supports Ubuntu/Debian only."
  exit 1
fi

echo "[1/6] System update"
sudo apt-get update -y
sudo apt-get upgrade -y

echo "[2/6] Install Docker + Compose plugin"
sudo apt-get install -y ca-certificates curl gnupg lsb-release
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list >/dev/null
sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

echo "[3/6] Enable Docker"
sudo systemctl enable docker
sudo systemctl start docker

echo "[4/6] Add current user to docker group"
sudo usermod -aG docker "$USER"

echo "[5/6] Open HTTP/HTTPS via UFW (if active)"
if command -v ufw >/dev/null 2>&1; then
  sudo ufw allow 80/tcp || true
  sudo ufw allow 443/tcp || true
fi

echo "[6/6] Done"
echo "Log out and log back in so docker group membership takes effect."
echo "Then deploy FotoBoo using deploy/oracle/README.md"
