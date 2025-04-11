# Imagem base para construção
FROM node:20-alpine AS builder

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY package.json package-lock.json* ./

# Instalar dependências
RUN npm ci

# Copiar código fonte
COPY . .

# Construir a aplicação
RUN npm run build

# Imagem final
FROM node:20-alpine AS runner

# Definir diretório de trabalho
WORKDIR /app

# Definir variáveis de ambiente
ENV NODE_ENV production

# Adicionar usuário não-root para segurança
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copiar arquivos necessários da imagem de construção
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# Definir usuário não-root
USER nextjs

# Expor porta
EXPOSE 3000

# Definir variável de ambiente para a porta
ENV PORT 3000

# Definir comando para iniciar a aplicação
CMD ["node", "server.js"]
