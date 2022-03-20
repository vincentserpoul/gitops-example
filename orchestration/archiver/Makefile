.PHONY: test bench bench-compare build docker-build lint changelog-gen  sec-scan upgrade db-init sql-gen

PROJECT_NAME?=orchestration
APP_NAME?=archiver

########
# test #
########

test:
	go test ./pkg/... -race -cover && \
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

######
# db #
######

db-init: 
	@(printf "Enter pass for db: "; read -rs DB_PASSWORD && \
	printf "\nEnter environment suffix(_dev, _local...): "; read DB_SUFFIX &&\
	sed \
	-e "s/DB_PASSWORD/$$DB_PASSWORD/g" \
	-e "s/DB_SUFFIX/$$DB_SUFFIX/g" \
	./sql/init/init.sql | \
	PGPASSWORD=$$DB_PASSWORD psql -h localhost -p 5436 -U postgres -f -)

########
# sqlc #
########

sql-gen:
	sqlc -f ./sql/sqlc.yaml generate