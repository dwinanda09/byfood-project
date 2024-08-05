.PHONY: build up down logs health test test-coverage clean db-backup db-reset

# Build all Docker images
build:
	docker-compose build

# Start all services in background
up:
	docker-compose up -d

# Stop services (preserve data)
down:
	docker-compose down

# Show logs for all services
logs:
	docker-compose logs -f

# Check service health status
health:
	@echo "=== Service Health Check ==="
	@docker-compose ps
	@echo "\n=== PostgreSQL Health ==="
	@docker-compose exec postgres pg_isready -U postgres || echo "PostgreSQL not ready"
	@echo "\n=== Backend Health ==="
	@curl -s http://localhost:8080/health || echo "Backend not responding"
	@echo "\n=== Frontend Health ==="
	@curl -s http://localhost:3000 > /dev/null && echo "Frontend: OK" || echo "Frontend not responding"

# Run backend unit tests
test:
	@echo "Running backend unit tests..."
	cd backend && go test ./... -v

# Run backend tests with coverage report
test-coverage:
	@echo "Running backend tests with coverage..."
	cd backend && go test ./... -v -coverprofile=coverage.out
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at backend/coverage.html"

# Create database backup
db-backup:
	@echo "Creating database backup..."
	docker-compose exec postgres pg_dump -U postgres byfood_library > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Backup created: backup_$(shell date +%Y%m%d_%H%M%S).sql"

# Reset database (with confirmation)
db-reset:
	@read -p "Are you sure you want to reset the database? This will delete all data! (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		docker-compose down; \
		docker volume rm byfood-project_postgres_data || true; \
		docker-compose up -d; \
		echo "Database reset complete"; \
	else \
		echo "Database reset cancelled"; \
	fi

# Remove all containers and data (with confirmation)
clean:
	@read -p "Are you sure you want to remove all containers and data? (y/N): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		docker-compose down -v; \
		docker-compose rm -f; \
		docker volume prune -f; \
		echo "Cleanup complete"; \
	else \
		echo "Cleanup cancelled"; \
	fi