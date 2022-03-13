.PHONY: test build docker-build sec-scan

PROJECT_NAME?=orchestration
APP_NAME?=orchestrator

########
# test #
########

test:
	go test ./pkg/... -race && \
	go test ./pkg/... -leak

bench:
	

bench-compare:
	benchstat 

#########
# build #
#########

build: test docker-build 
	@echo "tag your build with - git tag $(APP_NAME)-vX.X.X"

docker-build:
	docker pull gcr.io/distroless/static && \
	docker build -t $(APP_NAME) ./

#######
# sec #
#######

sec-scan:
	trivy fs ./

############
# upgrades #
############

upgrade:
	go mod tidy && \
	go get -t -u ./... && \
	go mod tidy
