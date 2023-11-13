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
test:
	@echo Conifig Environment Variable
	export CONFIG_PATH=./cmds/session-monitor/config.yaml 
	@echo Run tests
	go test -v  -covermode=count -coverprofile=coverage.out ./... 
	@echo Generate Html
	grep -v -E -f .covignore coverage.out > coverage.filtered.out
	mv coverage.filtered.out coverage.out
	go tool cover -html coverage.out -o coverage.html
	@echo Generate Xml
	gocover-cobertura < coverage.out > coverage.xml

rebuild-image: clean-image build-image

clean-image:
	docker image rm session-monitor-k8s

build-image:
	docker build -t session-monitor-k8s --file Dockerfile .