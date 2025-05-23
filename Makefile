MIGRATIONS_PATH = internal/db/migrations
DATABASE_URL = postgres://postgres:root@postgres:5432/go_social?sslmode=disable

migration:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

migrate_up:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) up

migrate_down:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) down $(filter-out $@,$(MAKECMDGOALS))

migrate_version:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) version

migrate_fix:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) force $(filter-out $@,$(MAKECMDGOALS))

%:
	@: