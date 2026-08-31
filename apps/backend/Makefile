# Makefile — apps/backend
# Jalankan `make help` untuk daftar perintah.

# Muat variabel dari .env kalau ada (dipakai target migrate-*).
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

BINARY      := questday
MAIN        := ./cmd/api
MIGRATIONS  := ./migrations

.PHONY: help
help: ## Tampilkan daftar perintah
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: dev
dev: ## Jalankan server dengan hot reload (air)
	air

.PHONY: run
run: ## Jalankan server sekali tanpa reload
	go run $(MAIN)

.PHONY: build
build: ## Compile binary ke ./bin
	go build -o bin/$(BINARY) $(MAIN)

.PHONY: tidy
tidy: ## Rapikan go.mod / go.sum
	go mod tidy

.PHONY: fmt
fmt: ## Format kode
	go fmt ./...

.PHONY: vet
vet: ## Analisis statis bawaan Go
	go vet ./...

.PHONY: test
test: ## Jalankan test
	go test ./... -race -cover

# --- Migrasi (butuh CLI golang-migrate terpasang) --------------------------
# Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: migrate-up
migrate-up: ## Terapkan semua migrasi
	migrate -path $(MIGRATIONS) -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## Rollback satu langkah
	migrate -path $(MIGRATIONS) -database "$(DATABASE_URL)" down 1

.PHONY: migrate-create
migrate-create: ## Buat file migrasi baru: make migrate-create name=nama_migrasi
	migrate create -ext sql -dir $(MIGRATIONS) -seq $(name)

.PHONY: migrate-force
migrate-force: ## Paksa versi (perbaiki state 'dirty'): make migrate-force version=N
	migrate -path $(MIGRATIONS) -database "$(DATABASE_URL)" force $(version)
