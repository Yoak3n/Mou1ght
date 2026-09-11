# 生产环境部署指南

本指南介绍如何将 Mou1ght 部署到带有域名和 SSL 证书的生产服务器上。

## 架构概览

```
                        ┌─────────────────────────┐
                        │     Nginx (Docker)      │
                  :80   │                         │  :443
                  ─────▶│  HTTP → 301 HTTPS       │◀────── 用户浏览器
                        │                         │
                        │  / ───────▶ client:3000 │
                        │  /api/ ──▶ backend:10420│
                        │  /admin/ ▶ backend:10420│
                        │  /upload/▶ backend:10420│
                        └─────────────────────────┘
                              │             │
                    ┌─────────┘             └─────────┐
                    ▼                                 ▼
            ┌──────────────┐                ┌──────────────────┐
            │    client    │                │     backend      │
            │  (Next.js)   │                │ (Go + Admin UI)  │
            │   port 3000  │                │    port 10420    │
            └──────────────┘                │                  │
                                            │  SQLite (volume) │
                                            └──────────────────┘
```

## 前置条件

- 一台 Linux 服务器（推荐 Ubuntu 22.04+）
- 已解析到服务器的域名
- SSH 访问权限（root 或 sudo）

## 第一步：安装 Docker

```bash
curl -fsSL https://get.docker.com | sh
sudo systemctl enable docker
sudo systemctl start docker

docker --version
docker compose version
```

## 第二步：获取代码

```bash
cd /opt
sudo git clone --recursive https://github.com/Yoak3n/Mou1ght.git
sudo chown -R $USER:$USER Mou1ght
cd Mou1ght
```

> `--recursive` 会自动拉取 `frontend` 子模块。

## 第三步：配置应用

编辑 `config.yaml`，根据需要修改站点信息、数据库配置、JWT 密钥等：

```bash
vi config.yaml
```

参考 `config.yaml` 中已有的注释字段进行配置。

## 第四步：配置 Nginx

### 4.1 修改域名

编辑 `nginx/conf.d/default.conf`，将所有 `your-domain.com` 替换为你的实际域名：

```bash
sed -i 's/your-domain.com/blog.example.com/g' nginx/conf.d/default.conf
```

### 4.2 准备证书目录

```bash
mkdir -p nginx/certs
```

## 第五步：配置域名与 SSL 证书

先替换 nginx 配置里的占位域名（`server_name` 与证书路径中的 `your-domain.com`）：

```bash
sed -i 's/your-domain.com/你的域名/g' nginx/conf.d/default.conf
```

### 方式 A：Let's Encrypt 自动证书（推荐，无需手动续期）

适用于已有域名解析且 80 端口可达的情况。首次签发用 **webroot 认证器**（经 nginx 验证），后续 `certbot renew` 才能自动走通。

**1. 用临时 80-only 配置占位（证书还没签发，nginx 此时不能加载 443）：**

```bash
cp nginx/conf.d/default.conf /tmp/default.conf.bak
cat > nginx/conf.d/default.conf << 'EOF'
server {
    listen 80;
    server_name 你的域名;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 200 "Mou1ght - setting up SSL...";
    }
}
EOF
docker compose -f docker-compose.prod.yaml up -d nginx
```

**2. 签发证书（webroot 认证器会存进续签配置）：**

```bash
docker compose -f docker-compose.prod.yaml run --rm certbot \
  certonly --webroot \
  -w /var/www/certbot \
  -d 你的域名 \
  --email your@email.com \
  --agree-tos --no-eff-email --non-interactive
```

**3. 恢复完整配置并启动全部服务：**

```bash
mv /tmp/default.conf.bak nginx/conf.d/default.conf
docker compose -f docker-compose.prod.yaml up -d
```

**4. 之后无需任何手动操作**：证书由 compose 里的 certbot 服务自动续签。

> **自动续签机制**：certbot 服务每 12 小时检查一次，在证书到期前 30 天自动续签（Let's Encrypt 证书有效期 90 天，等于到期前自动换新）；nginx 每 12 小时重载一次，续签后的新证书会自动生效，全程零停机、无需手工干预。
>
> 前提：证书来自 ACME 协议 CA（Let's Encrypt / ZeroSSL 等）。如果是云厂商的非 ACME 免费证书（如阿里云/腾讯云，通常一年期），certbot 无法自动续签，请改用方式 B 或用厂商提供的 CLI/API 自动化。

### 方式 B：使用已有证书文件

如果你的证书不是 ACME 协议（或想手动管理），把证书放到 nginx 配置引用的路径下：

```bash
mkdir -p nginx/certs/live/你的域名
cp /path/to/fullchain.pem nginx/certs/live/你的域名/fullchain.pem
cp /path/to/privkey.pem   nginx/certs/live/你的域名/privkey.pem
chmod 600 nginx/certs/live/你的域名/*.pem
```

> 方式 B 不会自动续期，证书到期前需要自己按厂商流程更换后替换这两个文件（nginx 每 12h 重载会自动加载新文件）。

## 第六步：启动所有服务

```bash
docker compose -f docker-compose.prod.yaml up -d --build
```

首次启动需要构建镜像，耗时较长。

## 第七步：验证部署

```bash
# 检查所有容器状态
docker compose -f docker-compose.prod.yaml ps

# 应该看到 4 个容器都在运行：
# - nginx
# - certbot
# - backend
# - client

# 检查 HTTPS
curl -I https://blog.example.com

# 检查 API
curl https://blog.example.com/api/v1/setting/blog
```

浏览器访问：
- 前端客户端：`https://blog.example.com`
- 管理后台：`https://blog.example.com/admin`
- 后端 API：`https://blog.example.com/api/v1/`

## 常用运维命令

```bash
# 查看日志
docker compose -f docker-compose.prod.yaml logs -f backend
docker compose -f docker-compose.prod.yaml logs -f client
docker compose -f docker-compose.prod.yaml logs -f nginx

# 重启单个服务
docker compose -f docker-compose.prod.yaml restart backend

# 更新代码后重新部署
git pull
docker compose -f docker-compose.prod.yaml up -d --build

# 仅重建后端
docker compose -f docker-compose.prod.yaml up -d --build backend

# 停止所有服务
docker compose -f docker-compose.prod.yaml down

# 停止并清除所有数据（慎用）
docker compose -f docker-compose.prod.yaml down -v

# 手动续期证书
docker compose -f docker-compose.prod.yaml run --rm certbot renew
docker compose -f docker-compose.prod.yaml restart nginx
```

## 数据备份

SQLite 数据库存储在 Docker volume `mou1ght-data` 中，备份方法：

```bash
# 导出 volume
docker run --rm -v mou1ght-data:/data -v $(pwd):/backup alpine \
  tar czf /backup/mou1ght-data-$(date +%Y%m%d).tar.gz -C /data .

# 恢复 volume
docker run --rm -v mou1ght-data:/data -v $(pwd):/backup alpine \
  tar xzf /backup/mou1ght-data-20250101.tar.gz -C /data
```

## 故障排查

### 502 Bad Gateway

Nginx 无法连接后端或客户端。检查：

```bash
docker compose -f docker-compose.prod.yaml ps
docker compose -f docker-compose.prod.yaml logs backend
docker compose -f docker-compose.prod.yaml logs client
```

### SSL 证书错误

```bash
# 检查证书文件是否存在（路径与 nginx 配置的 live 目录一致）
ls -la nginx/certs/live/你的域名/

# 检查证书是否过期
openssl x509 -in nginx/certs/live/你的域名/fullchain.pem -noout -dates

# 查看自动续签是否正常（每 12h 一次）
docker compose -f docker-compose.prod.yaml logs certbot

# 强制续期（手动触发）
docker compose -f docker-compose.prod.yaml run --rm certbot renew --force-renewal --webroot -w /var/www/certbot
```

### 客户端无法连接 API

检查 Nginx 配置中的 `proxy_pass` 是否指向正确的内部服务地址。Docker 内部网络使用服务名（`backend`、`client`）作为主机名。
