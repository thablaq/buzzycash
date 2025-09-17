# ==========================
# VARIABLES
# ==========================
ROOT_DIR := ./
DB_MIGRATIONS_DIR := $(ROOT_DIR)internal/db/migrations
APP_NAME := buzzycash
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_CONTAINER := $(APP_NAME)_container

# ==========================
# GO COMMANDS
# ==========================
tidy:
	cd $(ROOT_DIR) && go mod tidy

vet:
	cd $(ROOT_DIR) && go vet ./...

# --- Swagger ---
.PHONY: swagger
swagger:
	swag init -g cmd/main.go -o ./docs --parseDependency --parseInternal

.PHONY: run
run: swagger
	cd $(ROOT_DIR) && go run cmd/main.go

build:
	cd $(ROOT_DIR) && go build -o bin/$(APP_NAME) cmd/main.go

# ==========================
# DOCKER COMMANDS
# ==========================
docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker run --rm -it -p 8080:8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

docker-sh:
	docker exec -it $(DOCKER_CONTAINER) sh

docker-stop:
	docker stop $(DOCKER_CONTAINER) || true

# ==========================
# DATABASE MIGRATIONS
# ==========================
create-migration:
ifdef name
	migrate create -ext sql -dir $(DB_MIGRATIONS_DIR) -seq $(name)
else
	@echo "Please provide a migration name, e.g., make create-migration name=add_users_table"
endif

migrate-up:
	set -a; source .env; set +a; \
	migrate -database "$$DATABASE_URL" -path $(DB_MIGRATIONS_DIR) up

migrate-down:
	set -a; source .env; set +a; \
	migrate -database "$$DATABASE_URL" -path $(DB_MIGRATIONS_DIR) down 1

migrate-to:
ifdef version
	set -a; source .env; set +a; \
	migrate -database "$$DATABASE_URL" -path $(DB_MIGRATIONS_DIR) goto $(version)
else
	@echo "Please provide a version number, e.g., make migrate-to version=3"
endif

migrate-drop:
	set -a; source .env; set +a; \
	migrate -database "$$DATABASE_URL" -path $(DB_MIGRATIONS_DIR) drop -f

# ==========================
# DATABASE SCHEMA DUMP + MIGRATION
# ==========================
dump:
	@echo "Dumping schema from database..."
	set -a; source .env; set +a; \
	pg_dump --schema-only --no-owner --no-privileges -d "$$DATABASE_URL" > schema.sql
	@echo "Schema dump written to schema.sql ✅"

	@echo "Running migrate.sh to generate migration files..."
	chmod +x migrate.sh
	./migrate.sh
	@echo "Migration files generated in $(DB_MIGRATIONS_DIR) ✅"
