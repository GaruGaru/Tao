DOCKER_IMAGE=garugaru/tao
BUILD=$(shell git rev-parse --short HEAD)


.PHONY: lint
lint:
	revive -config lint/revive.toml -exclude ./vendor/... -formatter stylish ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: build
build:
	go build

.PHONY: test
test:
	go test -v ./...

.PHONY: deps
deps:
	dep ensure

.PHONY: docker_up
docker_up:
	docker-compose -f docker/docker-compose.yml up


.PHONY: docker_build_dev
docker_build_dev:
	docker build -t ${DOCKER_IMAGE}:${BUILD} -f docker/Dockerfile.dev .


.PHONY: docker_build_prod
docker_build_prod:
	docker build -t ${DOCKER_IMAGE}:${BUILD} -t ${DOCKER_IMAGE}:latest -f docker/Dockerfile.prod .


.PHONY: docker_push
docker_push: docker_build_prod
	docker push ${DOCKER_IMAGE}:${BUILD}
