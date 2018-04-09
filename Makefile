BINARY=tao
COMPOSE=docker/docker-compose.yml
COMPOSE_TEST=docker/docker-compose-test.yml

build:
	go build -o ${BINARY}

docker_up:
	docker-compose -f ${COMPOSE} up

docker_test:
	docker-compose -f ${COMPOSE} -f ${COMPOSE_TEST} up

