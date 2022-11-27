mqup:
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:latest

mqstop:
	docker stop rabbitmq && docker rm -f rabbitmq

database:
	docker run -d --rm \
			--name database \
 			-p 5432:5432 \
 			-e POSTGRES_USER=pg \
 			-e POSTGRES_PASSWORD=pass \
 			-e POSTGRES_DB=crud \
 			-v db:/var/lib/postgresql/data \
 			-d postgres:latest
dockerup:
	docker-compose up -d --build

dockerstop:
	docker-compose stop
