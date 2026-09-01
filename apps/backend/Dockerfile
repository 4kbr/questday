# Dockerfile produksi (multi-stage). Untuk dev pakai `make dev` + air, bukan ini.

# --- build ---
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/questday ./cmd/api

# --- runtime ---
FROM alpine:3.20
RUN adduser -D -u 10001 app
COPY --from=build /out/questday /usr/local/bin/questday
USER app
EXPOSE 8080
ENTRYPOINT ["questday"]
