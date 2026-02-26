# Oracle Always Free Deployment

Deploy FotoBoo on an Oracle Cloud Always Free VM with Docker Compose and automatic HTTPS (Caddy).

For the fastest zero-DNS launch, you can use `deploy/oracle/deploy-nipio.sh` to deploy on `https://<VM_IP>.nip.io` automatically.

## Why this option

- Always-on free VM (no sleep on idle)
- Persistent block storage
- Full control over your production stack

## 1) Create the VM in Oracle Cloud

Recommended settings:

- Shape: `VM.Standard.A1.Flex` (Always Free)
- OCPU/RAM: start with `1 OCPU / 6 GB`
- Image: Ubuntu 22.04 LTS
- Public IP: enabled
- Add your SSH key during instance creation

Open inbound ports in Oracle VCN security list or NSG:

- `22` (SSH)
- `80` (HTTP)
- `443` (HTTPS)

## 2) Prepare the server

SSH in:

```bash
ssh ubuntu@<your-vm-public-ip>
```

Run first boot setup:

```bash
sudo apt-get update -y
sudo apt-get install -y git
git clone https://github.com/rolniuq/fotoboo.git
cd fotoboo
chmod +x deploy/oracle/oracle-first-boot.sh
./deploy/oracle/oracle-first-boot.sh
```

Log out and back in so Docker group permissions apply.

## 3) Configure environment

Copy env template and edit:

```bash
cd ~/fotoboo
cp .env.example .env
nano .env
```

Set at least:

- `DOMAIN=your.domain.com`
- `LETSENCRYPT_EMAIL=you@your.domain.com`
- `BASE_URL=https://your.domain.com`
- `USE_MINIO=false` (recommended first)

Keep `APP_BIND_ADDRESS=127.0.0.1` so only Caddy exposes public ports.

### Fast path (no DNS setup)

If you want to go live immediately without buying/configuring a domain:

```bash
cd ~/fotoboo
chmod +x deploy/oracle/deploy-nipio.sh
./deploy/oracle/deploy-nipio.sh
```

This uses `https://<public-ip>.nip.io` and writes the needed values into `.env`.

## 4) Point your domain to the VM

In your DNS provider:

- Create `A` record for your domain to VM public IP

Wait for DNS propagation, then verify:

```bash
dig +short your.domain.com
```

## 5) Launch production stack

Use the `prod` profile to start app + Caddy:

```bash
docker compose --profile prod up -d --build
docker compose ps
```

Expected ports:

- `app` exposed on `127.0.0.1:8080` only
- `caddy` exposed on `0.0.0.0:80` and `0.0.0.0:443`

Health checks:

```bash
curl -f https://your.domain.com/health
```

## 6) Enable autostart on reboot

Compose services already use `restart: unless-stopped`, and Docker is enabled at boot by the setup script.

After reboot, verify:

```bash
docker compose ps
```

## Operations

Update to latest code:

```bash
cd ~/fotoboo
git pull origin master
docker compose --profile prod up -d --build
```

View logs:

```bash
docker compose logs -f app
docker compose logs -f caddy
```

Stop stack:

```bash
docker compose --profile prod down
```

## Backups

At minimum, back up the persistent Docker volume holding `/app/data`.

Quick snapshot backup example:

```bash
docker run --rm -v fotoboo_fotoboo_data:/data -v "$PWD":/backup alpine \
  tar -czf /backup/fotoboo-data-$(date +%F).tar.gz -C /data .
```

Restore example:

```bash
docker run --rm -v fotoboo_fotoboo_data:/data -v "$PWD":/backup alpine \
  sh -c "cd /data && tar -xzf /backup/fotoboo-data-YYYY-MM-DD.tar.gz"
```

## Notes

- Browsers require HTTPS for camera access on non-localhost domains.
- Start without MinIO for simplicity; add MinIO profile later if needed.
