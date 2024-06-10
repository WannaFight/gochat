include .env

## help: print this help message
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## guard-%: check if variable is set, e.g. "guard-TEST-VAR"
guard-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Variable \"$*\" not set"; \
		exit 1; \
	fi

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrate/new
db/migrate/new: guard-name
	@echo "Creating migration files for ${name}..."
	@migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrate/up
db/migrate/up:
	@echo "Running up migrations..."
	@migrate -path ./migrations -database ${DB_DSN} up

## db/migrations/force: force migration version
.PHONY: db/migrate/force
db/migrate/force: guard-version
	@echo "Forcing to migrate to version=${version}"
	@migrate -path ./migrations -database ${DB_DSN} force ${version}
