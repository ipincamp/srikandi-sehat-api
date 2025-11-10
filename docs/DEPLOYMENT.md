# Deployment Guide

This guide covers deploying the Srikandi Sehat GraphQL API to production environments.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Environment Setup](#environment-setup)
3. [Database Setup](#database-setup)
4. [Application Deployment](#application-deployment)
5. [Docker Deployment](#docker-deployment)
6. [Production Checklist](#production-checklist)
7. [Monitoring & Logging](#monitoring--logging)
8. [Scaling](#scaling)
9. [Backup & Recovery](#backup--recovery)
10. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### System Requirements

- **Operating System**: Linux (Ubuntu 20.04+ or CentOS 8+)
- **RAM**: Minimum 2GB (4GB+ recommended)
- **CPU**: 2+ cores recommended
- **Disk**: 20GB+ available space
- **Network**: Stable internet connection

### Required Software

- Go 1.25.4+
- PostgreSQL 16+
- Reverse proxy (Nginx or Caddy)
- SSL/TLS certificates (Let's Encrypt recommended)
- Email service (Mailgun account for production)

---

## Environment Setup

### 1. Create Production User

```bash
# Create dedicated user for the application
sudo useradd -m -s /bin/bash srikandi
sudo passwd srikandi

# Switch to the user
sudo su - srikandi
```

### 2. Clone Repository

```bash
cd /opt
git clone https://github.com/ipincamp/srikandi-sehat-graphql.git
cd srikandi-sehat-graphql
```

### 3. Configure Environment Variables

```bash
cp .env.example .env.production
nano .env.production
```

**Production Environment Configuration:**

```bash
# Application
APP_ENV=production
APP_PORT=8000
APP_HOST=0.0.0.0
APP_TIMEZONE=Asia/Jakarta

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=srikandi_sehat_prod
DB_USER=srikandi_db_user
DB_PASS=<STRONG_PASSWORD_HERE>
DB_SSL_MODE=require
DB_TIMEZONE=Asia/Jakarta

# PASETO Token (MUST be 32 bytes)
TOKEN_SYMMETRIC_KEY=<GENERATE_32_BYTE_KEY_HERE>
TOKEN_ISSUER=srikandi-sehat-api
TOKEN_ACCESS_DURATION=15m
TOKEN_REFRESH_DURATION=720h

# Email Service (Mailgun for production)
MAIL_DRIVER=mailgun
MAIL_FROM_ADDRESS=no-reply@yourdomain.com
MAIL_FROM_NAME=Srikandi Sehat
MAILGUN_DOMAIN=mg.yourdomain.com
MAILGUN_API_KEY=<YOUR_MAILGUN_API_KEY>
```

### 4. Generate Secure Keys

```bash
# Generate 32-byte symmetric key for PASETO
openssl rand -base64 32

# Generate strong database password
openssl rand -base64 24
```

---

## Database Setup

### 1. Install PostgreSQL

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install -y postgresql-16 postgresql-contrib
sudo systemctl enable postgresql
sudo systemctl start postgresql
```

**CentOS/RHEL:**
```bash
sudo dnf install -y postgresql16-server postgresql16-contrib
sudo postgresql-16-setup initdb
sudo systemctl enable postgresql-16
sudo systemctl start postgresql-16
```

### 2. Create Database and User

```bash
# Switch to postgres user
sudo -u postgres psql

# Create database user
CREATE USER srikandi_db_user WITH PASSWORD 'your_secure_password';

# Create database
CREATE DATABASE srikandi_sehat_prod;

# Grant privileges
GRANT ALL PRIVILEGES ON DATABASE srikandi_sehat_prod TO srikandi_db_user;

# Exit
\q
```

### 3. Configure PostgreSQL

Edit `postgresql.conf`:

```bash
sudo nano /etc/postgresql/16/main/postgresql.conf
```

**Recommended Settings:**

```ini
# Connection Settings
listen_addresses = 'localhost'
max_connections = 100

# Memory Settings
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
work_mem = 4MB

# WAL Settings
wal_buffers = 16MB
checkpoint_completion_target = 0.9

# Logging
logging_collector = on
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_rotation_age = 1d
log_rotation_size = 100MB
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
log_min_duration_statement = 1000
```

Edit `pg_hba.conf`:

```bash
sudo nano /etc/postgresql/16/main/pg_hba.conf
```

**Add:**

```
# TYPE  DATABASE                USER                ADDRESS         METHOD
local   srikandi_sehat_prod     srikandi_db_user                    md5
host    srikandi_sehat_prod     srikandi_db_user    127.0.0.1/32    md5
```

Restart PostgreSQL:

```bash
sudo systemctl restart postgresql
```

### 4. Run Migrations

```bash
cd /opt/srikandi-sehat-graphql
make migrate
```

---

## Application Deployment

### 1. Build the Application

```bash
cd /opt/srikandi-sehat-graphql
make build
```

### 2. Create Systemd Service

```bash
sudo nano /etc/systemd/system/srikandi-sehat.service
```

**Service Configuration:**

```ini
[Unit]
Description=Srikandi Sehat GraphQL API
After=network.target postgresql.service
Wants=postgresql.service

[Service]
Type=simple
User=srikandi
Group=srikandi
WorkingDirectory=/opt/srikandi-sehat-graphql
EnvironmentFile=/opt/srikandi-sehat-graphql/.env.production
ExecStart=/opt/srikandi-sehat-graphql/bin/srikandisehat
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=srikandi-sehat

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/srikandi-sehat-graphql

[Install]
WantedBy=multi-user.target
```

### 3. Enable and Start Service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service (start on boot)
sudo systemctl enable srikandi-sehat

# Start service
sudo systemctl start srikandi-sehat

# Check status
sudo systemctl status srikandi-sehat

# View logs
sudo journalctl -u srikandi-sehat -f
```

---

## Docker Deployment

### 1. Create Dockerfile

**File:** `Dockerfile`

```dockerfile
# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy schema files
COPY --from=builder /app/internal/adapters/driving/graphql/schema ./internal/adapters/driving/graphql/schema

EXPOSE 8000

CMD ["./main"]
```

### 2. Create Docker Compose for Production

**File:** `docker-compose.prod.yml`

```yaml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: srikandi-sehat-api
    restart: always
    ports:
      - "127.0.0.1:8000:8000"
    environment:
      - APP_ENV=production
      - APP_PORT=8000
      - DB_HOST=db
      - DB_PORT=5432
      - DB_NAME=${DB_NAME}
      - DB_USER=${DB_USER}
      - DB_PASS=${DB_PASS}
      - DB_SSL_MODE=disable
      - TOKEN_SYMMETRIC_KEY=${TOKEN_SYMMETRIC_KEY}
      - TOKEN_ISSUER=${TOKEN_ISSUER}
      - TOKEN_ACCESS_DURATION=${TOKEN_ACCESS_DURATION}
      - TOKEN_REFRESH_DURATION=${TOKEN_REFRESH_DURATION}
      - MAIL_DRIVER=${MAIL_DRIVER}
      - MAIL_FROM_ADDRESS=${MAIL_FROM_ADDRESS}
      - MAIL_FROM_NAME=${MAIL_FROM_NAME}
      - MAILGUN_DOMAIN=${MAILGUN_DOMAIN}
      - MAILGUN_API_KEY=${MAILGUN_API_KEY}
    depends_on:
      db:
        condition: service_healthy
    networks:
      - srikandi-network
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  db:
    image: postgres:16-alpine
    container_name: srikandi-sehat-db
    restart: always
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASS}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - srikandi-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
      interval: 10s
      timeout: 5s
      retries: 5
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

volumes:
  postgres_data:
    driver: local

networks:
  srikandi-network:
    driver: bridge
```

### 3. Deploy with Docker

```bash
# Build and start
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f app

# Stop
docker-compose -f docker-compose.prod.yml down

# Update and restart
git pull
docker-compose -f docker-compose.prod.yml up -d --build
```

---

## Reverse Proxy Setup

### Nginx Configuration

**File:** `/etc/nginx/sites-available/srikandi-sehat`

```nginx
upstream srikandi_backend {
    server 127.0.0.1:8000;
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name api.yourdomain.com;
    
    return 301 https://$server_name$request_uri;
}

# HTTPS server
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name api.yourdomain.com;

    # SSL certificates
    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;
    
    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Logging
    access_log /var/log/nginx/srikandi-access.log;
    error_log /var/log/nginx/srikandi-error.log;

    # GraphQL endpoint
    location /query {
        proxy_pass http://srikandi_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # GraphQL Playground (disable in production or add auth)
    location / {
        # Uncomment to disable playground in production
        # return 404;
        
        proxy_pass http://srikandi_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://srikandi_backend/query;
        access_log off;
    }
}
```

**Enable site:**

```bash
sudo ln -s /etc/nginx/sites-available/srikandi-sehat /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### SSL Certificate (Let's Encrypt)

```bash
# Install Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d api.yourdomain.com

# Auto-renewal (already set up by certbot)
sudo certbot renew --dry-run
```

---

## Production Checklist

### Security

- [ ] Use strong, unique passwords
- [ ] Generate secure 32-byte token symmetric key
- [ ] Enable SSL/TLS (HTTPS only)
- [ ] Configure firewall (allow only 80, 443, 22)
- [ ] Disable GraphQL Playground in production
- [ ] Set up rate limiting
- [ ] Enable CORS restrictions
- [ ] Keep dependencies updated
- [ ] Use environment variables for secrets
- [ ] Enable database SSL mode
- [ ] Configure security headers
- [ ] Set up fail2ban for SSH

### Performance

- [ ] Configure database connection pool
- [ ] Enable database query caching
- [ ] Set up CDN for static assets
- [ ] Configure Nginx caching
- [ ] Enable gzip compression
- [ ] Optimize database indexes
- [ ] Monitor query performance
- [ ] Set up DataLoader for N+1 prevention

### Monitoring

- [ ] Set up application logging
- [ ] Configure log rotation
- [ ] Set up error tracking (Sentry, etc.)
- [ ] Monitor database performance
- [ ] Set up uptime monitoring
- [ ] Configure alerting
- [ ] Monitor disk space
- [ ] Track API metrics

### Backup

- [ ] Set up automated database backups
- [ ] Test backup restoration
- [ ] Store backups off-site
- [ ] Document recovery procedures
- [ ] Version control for migrations

### Documentation

- [ ] Document deployment process
- [ ] Create runbook for common issues
- [ ] Document environment variables
- [ ] Maintain API changelog
- [ ] Document recovery procedures

---

## Monitoring & Logging

### Application Logs

```bash
# View systemd logs
sudo journalctl -u srikandi-sehat -f

# View Docker logs
docker-compose -f docker-compose.prod.yml logs -f app

# View Nginx logs
sudo tail -f /var/log/nginx/srikandi-access.log
sudo tail -f /var/log/nginx/srikandi-error.log
```

### Database Logs

```bash
# PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-16-main.log
```

### Log Rotation

**File:** `/etc/logrotate.d/srikandi-sehat`

```
/var/log/srikandi-sehat/*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 srikandi srikandi
    sharedscripts
    postrotate
        systemctl reload srikandi-sehat
    endscript
}
```

---

## Scaling

### Horizontal Scaling

1. **Load Balancer Setup:**

```nginx
upstream srikandi_cluster {
    least_conn;
    server 10.0.1.10:8000 weight=3;
    server 10.0.1.11:8000 weight=3;
    server 10.0.1.12:8000 weight=2;
}

server {
    location / {
        proxy_pass http://srikandi_cluster;
    }
}
```

2. **Database Read Replicas:**

Configure read replicas for read-heavy workloads and use primary for writes.

3. **Caching Layer:**

Add Redis for session storage and query caching.

### Vertical Scaling

- Increase server resources (CPU, RAM)
- Optimize database queries
- Tune PostgreSQL configuration
- Increase connection pool size

---

## Backup & Recovery

### Automated Database Backup

**File:** `/usr/local/bin/backup-db.sh`

```bash
#!/bin/bash

BACKUP_DIR="/backup/postgres"
DB_NAME="srikandi_sehat_prod"
DB_USER="srikandi_db_user"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/backup_$DATE.sql.gz"

# Create backup directory
mkdir -p $BACKUP_DIR

# Perform backup
pg_dump -U $DB_USER -d $DB_NAME | gzip > $BACKUP_FILE

# Remove backups older than 30 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete

echo "Backup completed: $BACKUP_FILE"
```

**Crontab:**

```bash
# Daily backup at 2 AM
0 2 * * * /usr/local/bin/backup-db.sh >> /var/log/backup.log 2>&1
```

### Restore Database

```bash
# Restore from backup
gunzip -c /backup/postgres/backup_20251110_020000.sql.gz | psql -U srikandi_db_user -d srikandi_sehat_prod
```

---

## Troubleshooting

### Application Won't Start

```bash
# Check logs
sudo journalctl -u srikandi-sehat -n 50

# Check port availability
sudo netstat -tulpn | grep 8000

# Check environment variables
sudo systemctl show srikandi-sehat -p Environment

# Test configuration
cd /opt/srikandi-sehat-graphql
./bin/srikandisehat
```

### Database Connection Issues

```bash
# Test database connection
psql -U srikandi_db_user -d srikandi_sehat_prod -h localhost

# Check PostgreSQL status
sudo systemctl status postgresql

# Check PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-16-main.log
```

### High Memory Usage

```bash
# Check memory usage
free -h
htop

# Check application memory
ps aux | grep srikandisehat

# Restart service
sudo systemctl restart srikandi-sehat
```

### Slow Queries

```bash
# Enable slow query logging in PostgreSQL
# Edit postgresql.conf:
log_min_duration_statement = 1000

# Check slow queries
sudo tail -f /var/log/postgresql/postgresql-16-main.log | grep "duration:"
```

---

## Support

For production support:
- Check documentation at `/docs`
- Review logs for error messages
- Contact development team
- Create issue on GitHub

---

## Changelog

### v1.0.0 (2025-11-10)
- Initial production release
- Complete deployment documentation
- Systemd service configuration
- Docker deployment support
- Nginx reverse proxy setup
- SSL/TLS configuration
- Backup and monitoring setup
