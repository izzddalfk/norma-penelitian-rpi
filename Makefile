test-go:
	docker compose -f ./deploy/local/test/go.docker-compose.yml down -v --remove-orphans
	docker compose -f ./deploy/local/test/go.docker-compose.yml up --exit-code-from rest_api

run-go:
	docker compose -f ./deploy/local/run/go.docker-compose.yml down -v --remove-orphans
	docker compose -f ./deploy/local/run/go.docker-compose.yml up --build

run-node:
	docker-compose -f ./node/docker-compose.yml down -v --remove-orphans
	docker-compose -f ./node/docker-compose.yml up