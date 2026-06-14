include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	@docker compose up -d tm-postgres

env-down:
	@docker compose down tm-postgres

env-cleanup:
	@read -p "Очистить все volume-файлы окружения? [y/N]: " ans; \
	if [ "$$ans"="y" ]; then \
	  docker compose down tm-postgres && \
	  sudo rm -rf out/pgdata && \
	  echo "Volume-файлы окружения очищены."; \
	else \
	  echo "Отмена"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутвует параметр seq. " && \
		echo "Пример: make migrate-create seq=create_user"; \
		exit 1; \
	fi; \
	docker compose run --rm tm-postgres-migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"

migrate:
	@if [ -z "$(action)" ]; then \
    	echo "Отсутвует параметр action. " && \
    	echo "Пример: make migrate action=up"; \
    	exit 1; \
    fi; \
	docker compose run --rm tm-postgres-migrate \
    -path /migrations \
    -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@tm-postgres:5432/${POSTGRES_DB}?sslmode=disable \
    $(action)

migrate-up:
	@make migrate action=up

migrate-down:
	@make migrate action=down