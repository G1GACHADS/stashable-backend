start:
	docker-compose build && docker-compose up --remove-orphans

stop:
	docker-compose down

docker-clean:
	docker stop $$(docker ps -q)
	docker rm $$(docker ps -a -f status=exited -q)
	docker system prune -f
	docker system prune -f --volumes

docker-clean-containers:
	docker rm -f $$(docker ps -a -q)
	docker rm $$(docker ps -a -f status=exited -q)

docker-clean-app-image:
	docker rmi -f backend-api-1

docker-clean-images:
	docker rmi -f $$(docker images -a -q)

docker-clean-volumes:
	docker volume rm $$(docker volume ls -q)