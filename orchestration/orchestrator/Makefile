.PHONY: test build docker-build sec-scan

PROJECT_NAME?=orchestration
APP_NAME?=orchestrator

########
# test #
########

test:
	go test ./..

bench:
	

bench-compare:
	benchstat 

#########
# build #
#########

build: docker-build

docker-build:
	docker pull gcr.io/distroless/static && \
	docker build -t $(APP_NAME) ./

#######
# sec #
#######

sec-scan:
	trivy fs ./