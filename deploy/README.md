# 部署指南

## 目录结构
```
deploy/
├── docker/             # Docker 相关配置
├── nginx/              # Nginx 配置
└── README.md           # 本文档
```

## 环境要求
- Docker 20.10+
- Docker Compose 2.0+
- 磁盘空间: 至少 5GB
- 内存: 推荐 4GB+

## 快速开始

### 1. 克隆仓库
```bash
git clone https://your-repository-url.git
cd galaxy-empire-manager
```

### 2. 准备配置文件
```bash
# 创建配置目录
mkdir -p deploy/nginx/ssl configs

# 复制示例配置
cp configs.example/* configs/
```

### 3. 构建并启动服务
```bash
docker-compose -f deploy/docker/docker-compose.yml up -d --build
```

### 4. 验证部署
```bash
curl http://localhost:8080/health
# 应返回 {"status":"ok"}
```

## 生产环境部署

### 1. 配置 Nginx
编辑 `deploy/nginx/conf.d/game-admin.conf`:
```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /etc/nginx/ssl/your-cert.pem;
    ssl_certificate_key /etc/nginx/ssl/your-key.pem;

    location / {
        proxy_pass http://app:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. 放置 SSL 证书
```bash
cp your-ssl-cert.pem deploy/nginx/ssl/
cp your-ssl-key.pem deploy/nginx/ssl/
```

### 3. 启动生产环境
```bash
docker-compose -f deploy/docker/docker-compose.yml up -d --build
```

## 常用命令

### 查看服务状态
```bash
docker-compose -f deploy/docker/docker-compose.yml ps
```

### 查看日志
```bash
# 查看应用日志
docker logs -f game-admin

# 查看数据库日志
docker logs -f mongodb

# 查看 Nginx 日志
docker logs -f nginx
```

### 停止服务
```bash
docker-compose -f deploy/docker/docker-compose.yml down
```

### 更新服务
```bash
# 1. 拉取最新代码
git pull

# 2. 重新构建并启动
docker-compose -f deploy/docker/docker-compose.yml up -d --build
```

## 数据备份

### 备份 MongoDB
```bash
docker exec mongodb sh -c 'mongodump --archive' > db-backup-$(date +%Y%m%d).archive
```

### 恢复 MongoDB
```bash
docker exec -i mongodb sh -c 'mongorestore --archive' < db-backup.archive
```

## 监控与维护

### 资源使用情况
```bash
docker stats
```

### 清理旧镜像
```bash
docker image prune -a
```

## 故障排除

### 常见问题
1. **端口冲突**  
   检查 80、443、8080 端口是否被占用

2. **证书问题**  
   确认 SSL 证书路径正确且权限适当

3. **数据库连接失败**  
   检查 MongoDB 容器是否正常运行

## 安全建议
1. 定期更新 Docker 和系统补丁
2. 使用非 root 用户运行容器
3. 限制数据库外部访问
4. 定期轮换 SSL 证书

## 联系方式
技术支持：devops@yourcompany.com  
项目维护：AaronChenH 

## Nginx 配置说明

1. 基本配置
```nginx
location /api {
    proxy_pass http://app:8080;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

2. 启用 Gzip 压缩
```nginx
gzip on;
gzip_types text/plain text/css application/json application/javascript;
```

3. 静态文件缓存
```nginx
location /static {
    expires 1y;
    add_header Cache-Control "public";
}
``` 