.PHONY: test bench bench-compare lint build docker-build sec-scan upgrade db-init sql-gen

PROJECT_NAME?=simple
APP_NAME?=gohttp

CAROOT=$(shell mkcert -CAROOT)
LAST_MAIN_COMMIT_HASH=$(shell git rev-parse --short HEAD)
LAST_MAIN_COMMIT_TIME=$(shell git log main -n1 --format='%cd' --date='iso-strict')

########
# test #
########

test: test-race test-leak

test-race:
	go test ./pkg/... -race -cover

test-leak:
	go test ./pkg/... -leak

bench:
	go test ./pkg/... -bench=. -benchmem | tee ./bench.txt

bench-compare:
	benchstat ./bench.txt

########
# lint #
########

lint:
	golangci-lint run ./... --config=./.golangci.toml

#########
# build #
#########

build: lint test bench docker-build 
	@printf "\nyou can now deploy to your env of choice:\ncd deploy\nENV=dev make deploy-latest\n"

docker-build:
	docker pull gcr.io/distroless/static:latest-amd64 && \
	docker build \
		-t $(APP_NAME) \
		--build-arg LAST_MAIN_COMMIT_HASH=$(LAST_MAIN_COMMIT_HASH) \
		--build-arg LAST_MAIN_COMMIT_TIME=$(LAST_MAIN_COMMIT_TIME) \
		./

#######
# sec #
#######

sec-scan:
	trivy fs ./ && \
	trivy image $(APP_NAME):latest

############
# upgrades #
############

upgrade:
	go mod tidy && \
	go get -t -u ./... && \
	go mod tidy


###########
# swagger #
###########

swagger-gen:
	swag init -d ./cmd/serve --parseDependency