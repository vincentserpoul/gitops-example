.PHONY: test bench bench-compare lint build docker-build sec-scan upgrade db-init sql-gen

PROJECT_NAME?=orchestration
APP_NAME?=archiver

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
	golangci-lint run ./...

#########
# build #
#########

build: lint test bench docker-build 
	@printf "\nyou can now create a release:\ncd deploy\nmake release\n"

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

######
# db #
######

db-pg-init: 
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

###########
# swagger #
###########

swagger-gen:
	swag init -d ./cmd/serve --parseDependency