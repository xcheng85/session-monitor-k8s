install-tools:
	@echo installing tools
	@echo done

generate:
	@echo running code generation
	@go generate ./...
	@echo done

format:
	gofmt -l -s -w .

build-amd64: 
	@env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o ./cmds/session-monitor ./cmds/session-monitor
run:
	cd ./cmds/session-monitor && CONFIG_PATH=./config.yaml ./session-monitor

rebuild-image: clean-image build-image

clean-image:
	docker image rm session-monitor-k8s

build-image:
	docker build -t session-monitor-k8s --file Dockerfile .