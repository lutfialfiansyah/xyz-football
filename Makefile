.PHONY: run migrate-up migrate-down clean help docker-up docker-up-d docker-down docker-destroy docker-logs

# Default target
help:
	@echo "Available commands:"
	@echo ""
	@echo "  Local:"
	@echo "  make run            - Run the API server locally"
	@echo "  make migrate-up     - Run database migrations UP"
	@echo "  make migrate-down   - Run database migrations DOWN"
	@echo "  make clean          - Clean up the binaries"
	@echo ""
	@echo "  Docker:"
	@echo "  make docker-up      - Build & start all containers"
	@echo "  make docker-up-d    - Build & start all containers (detached)"
	@echo "  make docker-down    - Stop containers (keep data)"
	@echo "  make docker-destroy - Stop containers & delete volumes (fresh DB)"
	@echo "  make docker-logs    - Tail container logs"

# ==================== Local ====================
run:
	go run cmd/api/main.go

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

clean:
	rm -f bin/api
	rm -f bin/migrate

# ==================== Docker ====================
docker-up:
	docker-compose up --build

docker-up-d:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-destroy:
	docker-compose down -v

docker-logs:
	docker-compose logs -f
