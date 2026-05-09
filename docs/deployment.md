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

### 4.2 准备 SSL 证书目录

```bash
mkdir -p nginx/ssl
```

## 第五步：申请 SSL 证书

### 方式 A：Let's Encrypt 自动证书（推荐）

适用于已有域名解析且 80 端口可达的情况。

**1. 创建临时 Nginx 配置（仅用于证书申请）：**

```bash
cat > nginx/conf.d/default.conf << 'EOF'
server {
    listen 80;
    server_name blog.example.com;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 200 "Mou1ght - setting up SSL...";
        add_header Content-Type text/plain;
    }
}
EOF
```

**2. 启动 Nginx 和 certbot 网络：**

```bash
docker compose -f docker-compose.prod.yaml up -d nginx
```

**3. 申请证书：**

```bash
docker compose -f docker-compose.prod.yaml run --rm certbot \
  certonly --webroot \
  -w /var/www/certbot \
  -d blog.example.com \
  --email your@email.com \
  --agree-tos \
  --no-eff-email
```

**4. 复制证书到 Nginx 目录：**

```bash
docker compose -f docker-compose.prod.yaml run --rm certbot \
  sh -c "cp /etc/letsencrypt/live/blog.example.com/fullchain.pem /etc/nginx/ssl/fullchain.pem && \
         cp /etc/letsencrypt/live/blog.example.com/privkey.pem /etc/nginx/ssl/privkey.pem"
```

**5. 恢复完整的 Nginx 配置：**

```bash
sed -i 's/your-domain.com/blog.example.com/g' nginx/conf.d/default.conf
```

**6. 重启 Nginx：**

```bash
docker compose -f docker-compose.prod.yaml restart nginx
```

> Let's Encrypt 证书有效期 90 天。`docker-compose.prod.yaml` 中的 certbot 服务会每 12 小时检查续期，无需手动操作。

### 方式 B：使用已有证书文件

将你的证书文件放入 `nginx/ssl/` 目录：

```bash
cp /path/to/your/fullchain.pem nginx/ssl/
cp /path/to/your/privkey.pem nginx/ssl/
chmod 600 nginx/ssl/*.pem
```

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
# 检查证书文件是否存在
ls -la nginx/ssl/

# 检查证书是否过期
openssl x509 -in nginx/ssl/fullchain.pem -noout -dates

# 强制续期
docker compose -f docker-compose.prod.yaml run --rm certbot renew --force-renewal
```

### 客户端无法连接 API

检查 Nginx 配置中的 `proxy_pass` 是否指向正确的内部服务地址。Docker 内部网络使用服务名（`backend`、`client`）作为主机名。
