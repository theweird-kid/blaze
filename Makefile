.PHONY: run-local seed docker-up docker-down

# Run locally (requires MONGO_URI set)
run-scheduler:
	DB_NAME=blaze_db WORKER_URL=http://localhost:8080 go run cmd/scheduler/main.go

run-worker:
	DB_NAME=blaze_db go run cmd/worker/main.go

# Seed a job
seed:
	DB_NAME=blaze_db go run cmd/seeder/main.go

seed-cron:
	DB_NAME=blaze_db go run cmd/seeder/main.go cron

# Docker Compose
docker-up:
	docker compose up --build

docker-down:
	docker compose down
