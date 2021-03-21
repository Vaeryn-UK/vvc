.PHONY: test
test: build
	docker-compose run --rm builder go test -v ./...

.PHONY: sample
sample: build
	docker-compose run --rm builder sh -c 'go run cmd/compiler/main.go programs/loops.vvb > /tmp/sample && hexdump -C /tmp/sample && go run cmd/execute/main.go /tmp/sample'

.PHONY: builder-sh
builder-sh: build
	docker-compose run --rm builder bash

.PHONY: build
build: .env .cache
	docker-compose build

.cache:
	mkdir .cache

.env:
	cp .env.dist .env
	sed -i s/UID=USER_ID/UID=$$(id -u)/ .env