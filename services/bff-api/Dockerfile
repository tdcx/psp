# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/bff ./cmd/bff

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /
COPY --from=build /out/bff /bff
ENV PORT=3000
USER nonroot
EXPOSE 3000
ENTRYPOINT ["/bff"]