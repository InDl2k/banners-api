build:
	docker-compose -f "docker-compose.yml" build
start:
	docker-compose -f "docker-compose.yml" up -d
down:
	docker-compose -f "docker-compose.yml" down --rmi local
stop:
	docker-compose -f docker-compose.yml stop
restart:
	docker-compose -f docker-compose.yml stop
	docker-compose -f docker-compose.yml up -d
test:
	cd test
	go test -v