.DEFAULT_GOAL := help
help: ## Displays all the available commands
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

compose: ## Start dependencies (postgres, redis)
	docker-compose up -d

apiserver: ## Starts local server
	air -c .air.toml

prod: ## Spins up the containers in production mode
	docker-compose build && docker-compose -f docker-compose.yml \
		-f build/docker-compose.yml up --remove-orphans

stop: ## Stops the containers
	docker-compose down

populate: ## runs db populate script
	go run cmd/populate/goroutine_limiter.go cmd/populate/util.go cmd/populate/main.go

db-migrate:
	./hack/db/migrate_db.sh

db-migrate-down:
	./hack/db/migrate_db_down.sh

db-reset:
	docker stop stashable-backend-db-1
	docker rm stashable-backend-db-1
	docker volume rm stashable-backend_postgres_data
	docker-compose up -d

db-refresh: db-migrate-down db-migrate populate ## runs db refresh script

docker-clean: ## Cleans all
	docker stop $$(docker ps -q)
	docker rm $$(docker ps -a -f status=exited -q)
	docker system prune -f
	docker system prune -f --volumes

docker-clean-containers: ## Removes all the containers
	docker rm -f $$(docker ps -a -q)
	docker rm $$(docker ps -a -f status=exited -q)

docker-clean-app-image: ## Removes the app image
	docker rmi -f backend-api-1

docker-clean-images: ## Removes all the images
	docker rmi -f $$(docker images -a -q)

docker-clean-volumes: ## Removes all the volumes
	docker volume rm $$(docker volume ls -q)