docker:
	docker build -t tracker . 

docker test exec: 
	# docker kill $(docker ps -q)
	# docker system prune  
	docker build -t tracker .
	docker run --name tracker tracker 
	docker exec -it tracker bash

docker test:
	docker build -t tracker .
	docker run --name tracker tracker 
	
docker clean: 
	docker system prune
