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

build: docker-build 
	@echo "\n tag your build with - git tag $(APP_NAME)-vX.X.X"

docker-build:
	docker pull gcr.io/distroless/static && \
	docker build -t $(APP_NAME) ./

#######
# sec #
#######

sec-scan:
	trivy fs ./