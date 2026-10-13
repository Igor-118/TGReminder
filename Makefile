.PHONY: local-env docker-env start_docker run stop lint test db-up db-down migrate-up migrate-down migrate-status run-local

# env.example без ключей, которые переопределены в .secrets (чтобы не было дублей в .env).
ENV_EXAMPLE = { \
	  if [ -f .secrets ]; then \
	    keys=$$(sed -n 's/^\([A-Za-z_][A-Za-z0-9_]*\)=.*/\1/p' .secrets | paste -sd'|' -); \
	    grep -vE "^($$keys)=" env.example || true; \
	  else \
	    cat env.example; \
	  fi; \
	}

docker-env:
	{ \
	  echo "# This .env file is generated automatically for DOCKER environment by Makefile."; \
	  echo "# Do not edit it directly; edit env.example / .secrets and Makefile instead."; \
	  echo; \
	  $(ENV_EXAMPLE); \
	  if [ -f .secrets ]; then \
	    echo; \
	    echo "# --- secrets from .secrets (not committed) ---"; \
	    cat .secrets; \
	  fi; \
	} > .env

local-env:
	{ \
	  echo "# This .env file is generated automatically for LOCAL environment by Makefile."; \
	  echo "# Do not edit it directly; edit env.example / .secrets and Makefile instead."; \
	  echo; \
	  $(ENV_EXAMPLE) | sed \
	    -e 's|^POSTGRES_HOST=.*|POSTGRES_HOST=127.0.0.1|' \
	    -e 's|^\(GOOSE_DBSTRING=.*@\)[^:/]*|\1127.0.0.1|'; \
	  if [ -f .secrets ]; then \
	    echo; \
	    echo "# --- secrets from .secrets (not committed) ---"; \
	    cat .secrets; \
	  fi; \
	} > .env

start_docker:
	docker compose --env-file .env -f deployments/docker-compose.yml up -d --build

run: docker-env start_docker

stop:
	docker compose --env-file .env -f deployments/docker-compose.yml down

lint:
	golangci-lint run

test:
	go test ./... -coverpkg=./... -coverprofile=coverage.out
	go tool cover -func=coverage.out | tail -1
	go tool cover -html=coverage.out -o coverage.html

# --- Локальный запуск: приложение на хосте, Postgres в докере ---

COMPOSE = docker compose --env-file .env -f deployments/docker-compose.yml
# goose без установки: скачается и закэшируется при первом запуске
GOOSE ?= go run github.com/pressly/goose/v3/cmd/goose@v3.26.0

db-up:
	$(COMPOSE) up -d db
	@echo "waiting for postgres..."
	@until $(COMPOSE) exec -T db pg_isready -q; do sleep 1; done

db-down:
	$(COMPOSE) stop db

# goose сам читает GOOSE_* из .env (сначала make local-env)
migrate-up:
	$(GOOSE) up

migrate-down:
	$(GOOSE) down

migrate-status:
	$(GOOSE) status

run-local: local-env db-up migrate-up
	go run ./cmd/bot
