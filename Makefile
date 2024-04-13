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
start-local:
	go build cmd/main.go
	main
test-local:
	cd test
	go test -v