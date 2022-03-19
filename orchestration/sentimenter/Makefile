.PHONY: test bench bench-compare build docker-build lint changelog-gen  sec-scan upgrade db-init sql-gen

PROJECT_NAME?=orchestration
APP_NAME?=sentimenter

########
# test #
########

test:
	go test ./pkg/... -race && \
	go test ./pkg/... -leak

bench:
	go test ./pkg/... -bench=.

bench-compare:
	benchstat 

#########
# build #
#########

build: lint test docker-build 
	@echo "tag your build with - git tag $(APP_NAME)-vX.X.X"

docker-build:
	docker pull gcr.io/distroless/static && \
	docker build -t $(APP_NAME) ./

########
# lint #
########

lint:
	golangci-lint run ./...

#############
# changelog #
#############

changelog-gen:
	git cliff -c ./cliff.toml -o ./CHANGELOG.md --include-path "**/$(APP_NAME)/*" --repository ../../

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
