.DEFAULT_GOAL := help

ENV_FILE := .env

-include $(ENV_FILE)
export

export PROJECT_ROOT=$(shell pwd)
DOCKER_USER := $(shell id -u):$(shell id -g)

.PHONY: help require-env env-up env-down env-cleanup \
	migrate-create migrate-up migrate-down migrate-action

help: ## Показать список команд
	@awk 'BEGIN{FS=":.*?## "} /^[a-zA-Z0-9_-]+:.*?## /{n++; t[n]=$$1; d[n]=$$2; if(length($$1)>w) w=length($$1)} END{for(i=1;i<=n;i++) printf "  \033[36m%-*s\033[0m %s\n", w, t[i], d[i]}' $(firstword $(MAKEFILE_LIST))

require-env:
	@test -f $(ENV_FILE) || { echo "$(ENV_FILE) не найден — выполните: cp .env.example .env"; exit 1; }

env-up: require-env ## Поднять контейнер с БД
	@docker compose up -d go-todolist-database

env-port-forwarder: ## Поднять проброс портов базы данных
	@docker compose up -d go-todolist-port-forwarder

env-port-close: ## Остановить проброс портов базы данных
	@docker compose down go-todolist-port-forwarder

env-down: require-env ## Остановить контейнер с БД
	@docker compose down go-todolist-database

env-cleanup: require-env ## Удалить контейнеры вместе с volumes
	@read -p "Очистить все volumes? [y/N]: " answer; \
	if [ "$$answer" = "y" ]; then \
	  docker compose down -v && \
	  echo "Очищено"; \
	else \
	  echo "Отмена"; \
	fi

migrate-create: require-env ## Создать файл миграции (seq=<имя>)
	@test -n "$(seq)" || { echo "Не задан параметр 'seq'"; exit 1; }
	@docker compose run --rm --user "$(DOCKER_USER)" go-todolist-db-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up: ## Применить все миграции
	@$(MAKE) migrate-action action=up

migrate-down: ## Откатить последнюю миграцию
	@$(MAKE) migrate-action action="down 1"

migrate-action: require-env ## Выполнить произвольную команду migrate (action=<...>)
	@test -n "$(action)" || { echo "Не задан параметр 'action'"; exit 1; }
	@docker compose run --rm go-todolist-db-migrate \
		-path /migrations \
		-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@go-todolist-database:5432/$(POSTGRES_DB)?sslmode=disable \
		$(action)
