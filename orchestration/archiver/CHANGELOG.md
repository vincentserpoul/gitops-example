# Changelog

All notable changes to this project will be documented in this file.

## [archiver-v1.1.15] - 2022-04-05

### Bug Fixes

- Remove wrong jaeger host

### Miscellaneous Tasks

- Up deps

### Performance

- Reduce sample param for traces
- Increase concurrent users

## [archiver-v1.1.14] - 2022-04-04

### Bug Fixes

- Sample trace parent

## [archiver-v1.1.12] - 2022-04-04

### Documentation

- Makefile printf correction

### Performance

- Always sample trace

## [archiver-v1.1.11] - 2022-04-04

### Features

- Prod deploy
- New manifest prod
- Add prod run bench

### Performance

- Add more replicas in prod

### Sec

- Dep upg

## [archiver-v1.1.10] - 2022-04-04

### Features

- Add commit time and commit hash to build
- Add port for crdb migration

### Sec

- Add image sec scan

## [archiver-v1.1.9] - 2022-03-29

### Debug

- Dockerfile with build args

## [archiver-v1.1.8] - 2022-03-29

### Debug

- Add dump of buildinfo

## [archiver-v1.1.7] - 2022-03-29

### Bug Fixes

- Better init sql

### Features

- Add build flags

### Performance

- Add upx and strip to reduce image size

### Styling

- Golangci-lint fix for go workspace

## [archiver-v1.1.5] - 2022-03-28

### Features

- Add sample rate in ingress and default backend

## [archiver-v1.1.3] - 2022-03-28

### Performance

- Use parent sampler to make sure we don't sample 100% all the time, but just when the serv is reached directly

## [archiver-v1.1.2] - 2022-03-26

### Documentation

- Update TODO in README.md

### Features

- Add crdb init

### Refactor

- Remove read/write users

## [archiver-v1.1.0] - 2022-03-24

### Bug Fixes

- Benchmark default port to 3003 default config
- Remove /docs from .dockerignore

### Documentation

- Update makefile info for deploy
- Add cockroachdb scripts

### Features

- Enable swagger

### Miscellaneous Tasks

- Go mod tidy
- Update deps

### Performance

- Move to pgx instead of sql
- Move to cockroachdb
- Bench
- Bench

### Refactor

- Merge db config, as it should be done with pgbouncer or pgcat
- Migration simplified

### Styling

- Remove bad linter comment

## [archiver-v1.0.20] - 2022-03-22

### Performance

- Use github.com/goccy/go-json instead of encoding json

### Security

- Remove bench from docker

## [archiver-v1.0.19] - 2022-03-22

### Bug Fixes

- Use the new sig for the gnup

### Features

- Add common errors
- Add benches to be able to compare later on

### Miscellaneous Tasks

- Replace interface{} by any
- Add apply-latest to be able to fine tune yamls

### Refactor

- Use recently created errors

### Styling

- Replace interface{} by any

### Testing

- Add -benchmem for banches
- Fix benchmarks
- Rename create handler test
- Update profiling

## [archiver-v1.0.18] - 2022-03-22

### Documentation

- Release doc

### Features

- New with http simplified

### Miscellaneous Tasks

- Update golangci-lint
- Update deps

### Testing

- K6 bench update

## [archiver-v1.0.17] - 2022-03-20

### Miscellaneous Tasks

- Update CHANGELOG.md for archiver:archiver-v1.0.16

## [archiver-v1.0.16] - 2022-03-20

### Miscellaneous Tasks

- Update CHANGELOG.md for archiver:archiver-v1.0.15
- Update commit changelog

## [archiver-v1.0.15] - 2022-03-20

### Documentation

- Update

### Miscellaneous Tasks

- Handle changelog commit properly

## [archiver-v1.0.14] - 2022-03-20

### Documentation

- Update

## [archiver-v1.0.13] - 2022-03-20

### Miscellaneous Tasks

- Makefile update

## [archiver-v1.0.12] - 2022-03-20

### Miscellaneous Tasks

- Reorder changelog, commit

## [archiver-v1.0.11] - 2022-03-20

### Bug Fixes

- Add chown

### Documentation

- Makefile cleanup

### Miscellaneous Tasks

- V1.0.11

## [archiver-v1.0.10] - 2022-03-20

### Bug Fixes

- Add chown

### Miscellaneous Tasks

- V1.0.9
- Better messaging after build

## [archiver-v1.0.9] - 2022-03-20

### Bug Fixes

- Add chown

### Miscellaneous Tasks

- Add release and release-tag
- Add a silencer for git-cliff

## [archiver-v1.0.8] - 2022-03-20

### Bug Fixes

- Add user non directive to avoid permission issues

### Release

- V1.0.7

## [archiver-v1.0.7] - 2022-03-20

### Documentation

- Add CHANGELOG

### Miscellaneous Tasks

- Release
- Update golangci-lint version
- Update golangci-lint

### Testing

- Add cover
- Split test leak and race in the makefile

### Security

- Add nonroot tag for final docker image

## [archiver-v1.0.6] - 2022-03-18

### Bug Fixes

- Ingress conf by host

## [archiver-v1.0.5] - 2022-03-18

### Bug Fixes

- Remove \n

## [archiver-v1.0.4] - 2022-03-18

### Bug Fixes

- Duplicate entry

## [archiver-v1.0.3] - 2022-03-18

### Bug Fixes

- Handle no change for migration

### Documentation

- Gen

### Features

- Add lint before build

### Miscellaneous Tasks

- Cleanup .PHONY

### Release

- Archiver-v1.0.1

## [archiver-v1.0.1] - 2022-03-16

### Bug Fixes

- Add fatal to indicate migrations failure
- Add migrations to docker build

### Features

- Add init container

## [archiver-v1.0.0] - 2022-03-16

### Documentation

- Add todo

### Features

- Archiver with db interaction added
- Added changelog generation and generated it

### Miscellaneous Tasks

- Update deps
- Go 1.18

<!-- generated by git-cliff -->
