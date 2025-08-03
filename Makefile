.PHONY: test-integration
run-db:
	docker-compose -f docker-compose.yml up -d --wait
run-migrations:
	go run ./migrations/auto-migrate.go
stop-db:
	docker-compose -f docker-compose.yml down
test-integration:
	docker-compose -f ./test/docker-compose.test.yml up -d --wait
	go test -v ./internal/repository -tags=integration -count=1
	docker-compose -f ./test/docker-compose.test.yml down

test-migrations:
	docker-compose -f ./test/docker-compose.test.yml up -d --wait
	@echo "Applying test migrations..."
	go run test/integration/migrations/main.go
	@docker-compose -f ./test/docker-compose.test.yml exec postgres psql -U test -d test -c "\dt"
	docker-compose -f ./test/docker-compose.test.yml down