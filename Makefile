APP_NAME := cent-user
REST_PORT ?= 8080

DB_URL ?= postgres://root:root@localhost:5432/cent_user?sslmode=disable
MIGRATION_PATH := ./migrations

build:
	@echo "Building $(APP_NAME)"
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/$(APP_NAME) .

rest:
	@echo "Running REST on $(APP_NAME)" 
	go run main.go rest -p $(REST_PORT)

test:
	@echo "Running tests on $(APP_NAME)"
	go test -v -cover ./...

cleanup:
	@echo "removing /bin"
	rm -rf bin/

migrateadd:
	@echo "Adding new migration file $(NAME)"
	migrate create -ext sql -dir $(MIGRATION_PATH) $(NAME)

migrateup:
	@echo "Execute migrate up"
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose up

migratedown:
	@echo "Execute migrate down"
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose down

migratereset:
	@echo "Execute reset all migrations"
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" -verbose drop -f
