SERVER_AIR_CONF='./config/.air.toml'

dev_server:
	@trap 'rm -rf ./tmp/server' EXIT; air -c ${SERVER_AIR_CONF}

format:
	@gofmt -l -s -w . && go mod tidy

prepare:
	@./scripts/prepare.bash

compose_dev:
	@if [ ! -f .env ]; then make prepare; fi
	@docker-compose -f ${PWD}/deployments/docker-compose.dev.yml --env-file ${PWD}/.env up
