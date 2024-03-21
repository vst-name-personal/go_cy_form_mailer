build:
	docker build --pull --rm -f "dockerfile" -t gocyformmailer:latest "."
	docker compose -f "deployments/docker-compose.yaml" up -d --build 