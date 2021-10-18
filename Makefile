build:
	docker build .
test:
	go test ./... 
run:
	docker-compose up
kill:
	docker-compose down
prune:
	docker system prune -a

# force do not cache
.PHONY: build test run kill
