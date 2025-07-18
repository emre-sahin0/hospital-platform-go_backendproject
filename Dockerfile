# ==================== BUILD STAGE ====================
# Go uygulamasını derlemek için kullanılan stage
FROM golang:1.24-alpine AS builder

# Çalışma dizinini ayarla
WORKDIR /app

# Gerekli paketleri yükle (git, gcc vb.)
RUN apk add --no-cache git gcc musl-dev

# Go module dosyalarını kopyala (dependency caching için)
COPY go.mod go.sum ./

# Dependencies'leri indir (cache layer olarak)
RUN go mod download

# Kaynak kodları kopyala
COPY . .

# Swagger docs'ları generate et
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Uygulamayı derle (statik binary)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hospital-platform .

# ==================== RUNTIME STAGE ====================
# Minimal image ile çalıştırma stage'i
FROM alpine:latest

# Güvenlik ve timezone için gerekli paketler
RUN apk --no-cache add ca-certificates tzdata

# Çalışma dizinini oluştur
WORKDIR /root/

# Build stage'den binary'yi kopyala
COPY --from=builder /app/hospital-platform .

# Swagger static files'ları kopyala (eğer varsa)
COPY --from=builder /app/docs ./docs

# Uygulama için port açılımı
EXPOSE 8080

# Health check ekle
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Uygulamayı çalıştır
CMD ["./hospital-platform"] 