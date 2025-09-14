# ---- build ----
FROM node:20-alpine AS build
WORKDIR /app
ENV CI=true
COPY services/bff-api/package*.json ./
RUN npm ci
COPY services/bff-api/ .
RUN npm run build

# ---- runtime ----
FROM node:20-alpine
WORKDIR /app
ENV NODE_ENV=production
# Instala sólo prod
COPY services/bff-api/package*.json ./
RUN npm ci --omit=dev
# Copia artefactos
COPY --from=build /app/dist ./dist
# Si necesitas archivos como .env.example, config, etc, añádelos aquí
# USER node  # (usa usuario no root si no necesitas puertos <1024 y tienes permisos de fs)
EXPOSE 3000
CMD ["node","dist/index.js"]
