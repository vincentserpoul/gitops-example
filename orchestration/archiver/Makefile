.PHONY: test bench bench-compare lint build docker-build sec-scan upgrade db-init sql-gen

PROJECT_NAME?=orchestration
APP_NAME?=archiver

CAROOT=$(shell mkcert -CAROOT)
LAST_MAIN_COMMIT_HASH=$(shell git rev-parse --short HEAD)
LAST_MAIN_COMMIT_TIME=$(shell TZ=UTC git log main -n1 --format='%cd' --date='format-local:%Y-%m-%dT%H:%M:%S')

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
	docker pull gcr.io/distroless/static && \
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

######
# db #
######

db-pg-init: 
	@( \
	printf "Enter pass for db: "; read -rs DB_PASSWORD && \
	printf "\nEnter environment suffix(_dev, _local...): "; read DB_SUFFIX &&\
	sed \
	-e "s/DB_PASSWORD/$$DB_PASSWORD/g" \
	-e "s/DB_SUFFIX/$$DB_SUFFIX/g" \
	./sql/init/init.sql | \
	PGPASSWORD=$$DB_PASSWORD psql -h localhost -p 5436 -U postgres -f - \
	)

db-cockroachdb-rootkey:
	mkdir ./sql/crdb-certs && \
	kubectl cp cockroachdb/cockroachdb-0:cockroach-certs/ca.crt ./sql/crdb-certs/ca.crt -c db && \
	cockroach cert create-client \
		--certs-dir=./sql/crdb-certs \
		--ca-key=$(CAROOT)/rootCA-key.pem root

db-cockroachdb-init:
	@( \
	printf "Enter pass for db: "; read -rs DB_PASSWORD && \
	printf "\nEnter environment suffix(_dev, _local...): "; read DB_SUFFIX &&\
	printf "\nEnter port(26257...): "; read DB_PORT &&\
	sed \
	-e "s/DB_PASSWORD/$$DB_PASSWORD/g" \
	-e "s/DB_SUFFIX/$$DB_SUFFIX/g" \
	./sql/init/init.sql > ./sql/crdb-certs/init.sed.sql && \
	cockroach sql --certs-dir=./sql/crdb-certs -f ./sql/crdb-certs/init.sed.sql -p $$DB_PORT && \
	rm ./sql/crdb-certs/init.sed.sql \
	)

########
# sqlc #
########

sql-gen:
	sqlc -f ./sql/sqlc.yaml generate

###########
# swagger #
###########

swagger-gen:
	swag init -d ./cmd/serve --parseDependency