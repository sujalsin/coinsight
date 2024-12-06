.PHONY: up down build logs ps clean

# Start all services
up:
	docker-compose up -d

# Start all services and show logs
up-logs:
	docker-compose up

# Stop all services
down:
	docker-compose down

# Rebuild all services
build:
	docker-compose build

# Show logs for all services
logs:
	docker-compose logs -f

# Show running containers
ps:
	docker-compose ps

# Clean up all containers, volumes, and images
clean:
	docker-compose down -v --rmi all

# Start individual services
price-service:
	docker-compose up -d price-service

api-service:
	docker-compose up -d api-service

web-dashboard:
	docker-compose up -d web-dashboard

recommendation-engine:
	docker-compose up -d recommendation-engine

# Database management
db-reset:
	docker-compose exec api-service rails db:drop db:create db:migrate db:seed

# Run tests
test-all:
	docker-compose exec price-service go test ./...
	docker-compose exec api-service rails test
	docker-compose exec web-dashboard npm test
	docker-compose exec recommendation-engine go test ./...
