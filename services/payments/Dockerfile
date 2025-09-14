# ---- build ----
FROM golang:1.22 AS build
WORKDIR /src
# Aprovecha caché de dependencias
COPY services/payments/go.mod services/payments/go.sum ./
RUN go mod download
# Copia el código
COPY services/payments/ .
# Compila estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /out/app ./cmd/server

# ---- runtime ----
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /out/app /app/app
# Puerto (ajústalo a tu server)
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
