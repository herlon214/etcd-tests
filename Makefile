IMAGE=staging.repo.rcplatform.io/reynencourt/rc-discovery

run:
	touch /tmp/hosts
	go run ./cmd/rc-discovery/main.go -file /tmp/hosts

build: clean-dist build-linux build-macos

build-linux:
	GO111MODULE=on GOOS=linux go build  ./cmd/rc-discovery/main.go;mv main ./dist/linux/bin/rc-discovery

build-macos:
	GO111MODULE=on GOOS=darwin go build ./cmd/rc-discovery/main.go;mv main ./dist/linux/bin/rc-discovery

clean-dist:
	if [ -e ./dist ]; then rm -rf ./dist; fi; mkdir ./dist; mkdir -p ./dist/darwin/bin; mkdir -p ./dist/linux/bin

test:
	@go test -race -coverpkg=./pkg/... -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@rm coverage.out

release: docker-build
	docker push "${IMAGE}:latest"
	docker tag "${IMAGE}:latest" "${IMAGE}:`cat ./.version`"
	docker push "${IMAGE}:`cat ./.version`"

docker-build:
	docker build -t "${IMAGE}:latest" . --no-cache

docker-test: build-docker
	docker-compose up -d
	docker run --network host -v $PWD/testdata:/app/testdata --env-file=./.env ${IMAGE}

etcd-status:
	@go run github.com/etcd-io/etcd/etcdctl endpoint health