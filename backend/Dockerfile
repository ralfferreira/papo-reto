FROM golang:1.24-alpine AS builder

# Instalar dependências de build
RUN apk add --no-cache gcc musl-dev

# Configurar diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o papo-reto-api ./cmd/api

# Imagem final
FROM alpine:3.17

# Instalar dependências de runtime
RUN apk add --no-cache ca-certificates tzdata

# Configurar timezone
ENV TZ=America/Sao_Paulo

# Copiar binário compilado da imagem de build
COPY --from=builder /app/papo-reto-api /usr/local/bin/

# Copiar arquivo .env
COPY .env /

# Expor porta
EXPOSE 8080

# Comando para executar a aplicação
CMD ["papo-reto-api"]
