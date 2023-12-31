# session-monitor-k8s

## Pkgs
1. gin
2. uber/dig
3. config management: viper
4. cli: Cobra

## Bootstrap
```shell
go mod init github.com/xcheng85/session-monitor-k8s

mkdir -p docs docker internal cmds  

# create internal modules
cd internal && mkdir -p config logger module http test worker && cd ../

# shared modules
cd internal && mkdir -p config logger module http && cd ../

# config module (ultimate config management)
cd internal/config && touch config.go viper.go viper_test.go && cd ../../

# http module (http error management)
cd internal/http && touch utils.go utils_test.go && cd ../../

# logger module (centralized structural logging)
cd internal/logger && touch logger.go zap.go zap_test.go && cd ../../

# module (domain driven interface)
cd internal/module && touch module.go module_test.go && cd ../../

# test (helper for testing purpose only)
cd internal/test && touch utils.go utils_test.go && cd ../../

# worker (domain module coordinator)
cd internal/worker && touch syncer.go syncer_test.go && cd ../../

# domain driven modules
mkdir -p k8s pod node session

# k8s modules (deployment health checks)
cd k8s && mkdir -p internal && touch module.go module_test.go && cd ../
cd k8s/internal && mkdir -p handler rest && cd handler && touch k8s.go k8s_test.go && cd ../ && cd rest && touch router.go router_test.go && cd ../../

cd ../cmds
mkdir -p session-monitor
cd session-monitor
touch main.go app.go config.yaml dummy.yaml
```

## Build

```shell
make build-amd64
```

## Run
```shell
# in terminal #1
export CONTAINER_ID=
docker stop $CONTAINER_ID
docker rm $CONTAINER_ID
docker run --name=redis -p 6379:6379 redis:6.2.7

# in terminal #2
make run
```

## Testing
```shell
export CONFIG_PATH=./cmds/session-monitor/config.yaml
go test -v  -covermode=count -coverprofile=coverage.out ./...
grep -v -E -f .covignore coverage.out > coverage.filtered.out
mv coverage.filtered.out coverage.out
go tool cover -html coverage.out -o coverage.html
gocover-cobertura < coverage.out > coverage.xml
```