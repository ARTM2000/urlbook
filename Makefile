SERVER_AIR_CONF='./config/.air.toml'

dev_server:
	@trap 'rm -rf ./tmp/server' EXIT; air -c ${SERVER_AIR_CONF}

format:
	@gofmt -l -s -w . && go mod tidy

prepare:
	@./scripts/prepare.bash
