.PHONY: build
build: build-builder
	docker-compose run --rm builder go run cmd/vvc/main.go

.PHONY: builder-sh
builder-sh: build-builder
	docker-compose run --rm builder bash

.PHONY: build-builder
build-builder: .env .cache
	docker-compose build

.cache:
	mkdir .cache

.env:
	cp .env.dist .env
	sed -i s/UID=USER_ID/UID=$$(id -u)/ .env