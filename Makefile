.PHONY: db_run db_run_local local-env docker-env stop

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
	  echo "# Do not edit it directly; edit .env.example / .secrets and Makefile instead."; \
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
	    -e 's|^POSTGRES_HOST=.*|POSTGRES_HOST=127.0.0.1|'; \
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
