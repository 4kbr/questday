# Makefile — root monorepo QuestDay
# Satu pintu untuk infra dev (docker compose) + menjalankan app.
# Jalankan `make help` untuk daftar perintah.

# Muat .env.docker kalau ada (dipakai target compose & psql).
ifneq (,$(wildcard ./.env.docker))
	include .env.docker
	export
endif

# Default aman kalau .env.docker belum dibuat.
POSTGRES_USER ?= questday
POSTGRES_DB   ?= questday

COMPOSE := docker compose --env-file .env.docker -f docker-compose.dev.yml

.DEFAULT_GOAL := help

.PHONY: help
help: ## Tampilkan daftar perintah
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

# --- Infra dev (Postgres via docker compose) ------------------------------

# Auto-buat .env.docker dari template supaya --env-file tak gagal.
.env.docker: .env.docker.example
	@cp $< $@
	@echo "$@ dibuat dari template — cek isinya sebelum lanjut."

.PHONY: env-docker
env-docker: .env.docker ## Buat .env.docker dari template (kalau belum ada)

.PHONY: up
up: .env.docker ## Nyalakan infra dev (Postgres) di background
	$(COMPOSE) up -d

.PHONY: down
down: .env.docker ## Matikan infra dev (data tetap)
	$(COMPOSE) down

.PHONY: reset
reset: .env.docker ## Matikan infra dev + hapus volume (reset DB)
	$(COMPOSE) down -v

.PHONY: logs
logs: .env.docker ## Ikuti log infra dev
	$(COMPOSE) logs -f

.PHONY: ps
ps: .env.docker ## Status service infra dev
	$(COMPOSE) ps

.PHONY: psql
psql: .env.docker ## Masuk psql shell di container Postgres
	$(COMPOSE) exec postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

# --- Jalankan app --------------------------------------------------------

.PHONY: backend
backend: ## Jalankan backend (hot reload, :8080) — 'make up' dulu
	$(MAKE) -C apps/backend dev

.PHONY: frontend
frontend: ## Jalankan frontend (Vite dev, :5173)
	cd apps/frontend && npm run dev

.PHONY: dev
dev: ## Jalankan backend + frontend bersamaan (Ctrl-C hentikan dua-duanya)
	$(MAKE) -j2 backend frontend
