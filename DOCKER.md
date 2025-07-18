# 🐳 **Docker Deployment Guide**

Bu rehber, Hastane Yönetim Platformu'nu Docker ile nasıl çalıştıracağınızı anlatır.

---

## 🚀 **Hızlı Başlangıç**

### **1️⃣ Tüm Servisleri Çalıştır**
```bash
# Tüm servisleri başlat (Go app + PostgreSQL + Redis)
docker-compose up -d

# Logları izle
docker-compose logs -f

# Servislerin durumunu kontrol et
docker-compose ps
```

### **2️⃣ Uygulamaya Erişim**
```
🌐 Hospital API:    http://localhost:8080
📚 Swagger Docs:    http://localhost:8080/swagger/
🗄️ Database Admin:  http://localhost:8081 (optional)
🔴 Redis Admin:     http://localhost:8082 (optional)
```

### **3️⃣ Servisleri Durdur**
```bash
# Servisleri durdur
docker-compose down

# Servisleri durdur + volume'ları sil
docker-compose down -v
```

---

## 🔧 **Detaylı Kurulum**

### **📋 Ön Gereksinimler**
- **Docker:** 20.10+ 
- **Docker Compose:** 2.0+
- **RAM:** En az 2GB (tüm servisler için)
- **Disk:** En az 1GB boş alan

### **🔄 Adım Adım Setup**

#### **1. Repository'yi Klonla**
```bash
git clone https://github.com/vatansoft/hospital-platform.git
cd hospital-platform
```

#### **2. Environment Variables Kontrol Et**
```bash
# docker-compose.yml dosyasında environment değişkenleri:
# - DB_HOST=postgres
# - DB_PORT=5432
# - REDIS_HOST=redis
# - REDIS_PORT=6379
# - JWT_SECRET=super_secret_jwt_key_for_hospital_platform_2024
```

#### **3. Servisleri Başlat**
```bash
# Sadece core servisler (app + db + redis)
docker-compose up -d

# Admin tools ile birlikte (adminer + redis-commander)
docker-compose --profile tools up -d
```

#### **4. Database Migration Kontrol**
```bash
# Uygulama loglarını kontrol et
docker-compose logs app

# Başarılı migration mesajları:
# "Migration tamamlandı!"
# "Master data seeding completed!"
```

---

## 🐙 **Docker Compose Servisleri**

### **🏥 Core Services**

| Service | Port | Description |
|---------|------|-------------|
| **app** | 8080 | Go uygulaması (Hospital Platform) |
| **postgres** | 5432 | PostgreSQL 15 database |
| **redis** | 6379 | Redis cache server |

### **🛠️ Management Tools (Optional)**

| Service | Port | Description | Profile |
|---------|------|-------------|---------|
| **adminer** | 8081 | PostgreSQL web admin | `tools` |
| **redis-commander** | 8082 | Redis web admin | `tools` |

### **💾 Persistent Volumes**

| Volume | Description |
|--------|-------------|
| `postgres_data` | PostgreSQL data files |
| `redis_data` | Redis persistence files |

---

## 🔍 **Troubleshooting**

### **❗ Yaygın Problemler**

#### **1. Port Çakışması**
```bash
# Error: Port already in use
# Çözüm: docker-compose.yml'de portları değiştir

services:
  app:
    ports:
      - "8090:8080"  # 8080 yerine 8090 kullan
```

#### **2. Database Connection Error**
```bash
# Logs'ları kontrol et
docker-compose logs postgres
docker-compose logs app

# Database'in hazır olup olmadığını kontrol et
docker-compose exec postgres pg_isready -U postgres
```

#### **3. Redis Connection Error**
```bash
# Redis durumunu kontrol et
docker-compose exec redis redis-cli ping

# Yanıt: PONG (Redis çalışıyor)
```

#### **4. Build Hataları**
```bash
# Cache'i temizle ve yeniden build et
docker-compose build --no-cache app
docker-compose up -d
```

### **🔧 Debug Commands**

```bash
# Container'lara bağlan
docker-compose exec app sh          # Go uygulaması
docker-compose exec postgres psql -U postgres hospital_platform  # PostgreSQL
docker-compose exec redis redis-cli # Redis

# Logları detaylı izle
docker-compose logs -f --tail=100 app

# Resource kullanımını kontrol et
docker stats

# Network bağlantılarını kontrol et
docker network ls
docker network inspect hospital-network
```

---

## 🚀 **Production Deployment**

### **🔒 Security Checklist**

- [ ] **JWT Secret:** Güçlü, unique key kullan
- [ ] **Database Password:** Default password'u değiştir  
- [ ] **Redis Password:** Prodüksiyonda mutlaka password ekle
- [ ] **CORS:** Allowed origins'ı production domain ile sınırla
- [ ] **SSL/TLS:** Reverse proxy (nginx) ile HTTPS ekle

### **📈 Performance Optimization**

#### **1. Resource Limits**
```yaml
# docker-compose.yml'e ekle
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
```

#### **2. Health Checks**
```bash
# Health check endpoint'leri:
curl http://localhost:8080/health          # App health
curl http://localhost:8080/swagger/        # Swagger accessibility
```

#### **3. Backup Strategy**
```bash
# Database backup
docker-compose exec postgres pg_dump -U postgres hospital_platform > backup.sql

# Redis backup  
docker-compose exec redis redis-cli BGSAVE
```

### **🔄 CI/CD Integration**

```yaml
# .github/workflows/docker.yml örneği
name: Docker Build & Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build & Deploy
        run: |
          docker-compose build
          docker-compose up -d
          
      - name: Health Check
        run: |
          sleep 30
          curl -f http://localhost:8080/swagger/ || exit 1
```

---

## 📊 **Monitoring & Logs**

### **📈 Metrics Collection**

```bash
# Resource monitoring
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}"

# Disk usage
docker system df

# Log sizes
docker-compose logs --tail=0 -f | wc -l
```

### **🔍 Log Management**

```bash
# Structured logging
docker-compose logs --json app

# Log rotation (production)
# docker-compose.yml'e ekle:
services:
  app:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

---

## 🆘 **Support & Help**

### **📞 İletişim**
- **GitHub Issues:** https://github.com/vatansoft/hospital-platform/issues
- **Email:** info@vatansoft.com
- **Documentation:** http://localhost:8080/swagger/

### **🔗 Useful Links**
- **Docker Docs:** https://docs.docker.com/
- **Docker Compose:** https://docs.docker.com/compose/
- **PostgreSQL Docker:** https://hub.docker.com/_/postgres
- **Redis Docker:** https://hub.docker.com/_/redis

---

*Son Güncelleme: 2024 - VatanSoft Docker Guide v1.0* 