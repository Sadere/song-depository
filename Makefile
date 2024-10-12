DOCKER_COMPOSE_FILES ?= $(shell find docker -maxdepth 1 -type f -name "*.yaml" -exec printf -- '-f %s ' {} +; echo)

# Build and run docker containers
.PHONY: up
up:
	docker compose ${DOCKER_COMPOSE_FILES} up --build --detach

# Stop docker containers
.PHONY: down
down:
	docker compose ${DOCKER_COMPOSE_FILES} down