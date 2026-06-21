include .env

MIGRATIONS_PATH= /Users/huskiepuppy05/Development/Personal/go-job-application-tracker/cmd/migrations


.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))


.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable up


.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-force
migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable force 14 $(filter-out $@,$(MAKECMDGOALS))

# .PHONY: seed
# seed:
# 	@direnv allow /Users/ayomideademilua/Development/go_crash_course/go_social/.envrc
# 	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt


.PHONY: test
test:
	@go test -v ./...


.PHONY: clear-cache
clear-cache:
	@go clean -cache -testcache -modcache

.PHONY: test
test:
	@go test -v ./...


.PHONY: start-docker
start-docker:
	@docker compose up -d

.PHONY: stop-docker
stop-docker:
	@docker compose down

.PHONY: restart-docker
restart-docker:
	@docker compose down
	@docker compose up -d


.PHONY: clear-docker
clear-docker:
	@docker compose down
	@docker system prune -a
	@docker volume prune
	@docker volume prune -a
	@docker network prune
	@docker system prune -a --volumes
	@docker system prune -a --volumes