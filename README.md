# PSP MVP Bootstrap (SIF + Payments)

Este paquete es un **monorepo pnpm** con los servicios básicos para empezar:
- `billing-sif`: Diario SIF (append-only) con hash encadenado + QR y persistencia Postgres.
- `payments`: `payment_intents` (mock 3DS2) con persistencia Postgres.
- `bff-api`: BFF mínimo con health y proxy de desarrollo (opcional) hacia servicios.
- `webhooks`: receptor de webhooks con firma HMAC (solo skeleton).

Infra local:
- **Postgres 15** y **Redpanda** (Kafka) vía Docker Compose.

> Objetivo de esta parte: poder levantar local, emitir una SIF (`POST /v1/invoices`) y crear/confirmar un intent de pago.

## Requisitos
- Node.js 20
- pnpm 9
- Docker + Docker Compose

## Quickstart
```bash
# 1) Levanta Postgres y Redpanda
docker compose -f deploy/templates/docker-compose.local.yml up -d

# 2) Instala dependencias del monorepo
pnpm install

# 3) Arranca servicios en paralelo
pnpm dev
```

## Endpoints de prueba
```bash
# Healthchecks
curl -s localhost:3001/health | jq .           # billing-sif
curl -s localhost:3002/health | jq .           # payments

# Crear SIF
curl -s localhost:3001/v1/invoices  -H 'Content-Type: application/json'  -H 'Idempotency-Key: demo-1'  -d '{"merchant_id":"m1","series":"A","number":"0001","amount_cents":1234,"currency":"EUR"}' | jq .

# Crear y confirmar un intent
PI=$(curl -s localhost:3002/v1/payment_intents -H 'Content-Type: application/json'  -d '{"merchant_id":"m1","amount_cents":1999,"currency":"EUR"}' | jq -r .id)
curl -s -X POST localhost:3002/v1/payment_intents/$PI/confirm | jq .
```

## Notas
- Las tablas se **autocrean** al arrancar cada servicio (migración mínima).
- Redpanda/Kafka se incluye para el futuro **outbox**, pero no se usa todavía.
- Variables de entorno por servicio: ver `services/*/.env.example`.