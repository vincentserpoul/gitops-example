.PHONY: test build docker-build sec-scan

PROJECT_NAME?=choreography
APP_NAME?=hasher

########
# test #
########

test:
	go test ./..

#########
# build #
#########

build: docker-build

docker-build:
	docker build -t $(APP_NAME) ./

#######
# sec #
#######

sec-scan:
	trivy fs ./