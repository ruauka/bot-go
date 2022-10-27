mqup:
	docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.10-management

mqstop:
	docker stop rabbitmq && docker rm -f rabbitmq
