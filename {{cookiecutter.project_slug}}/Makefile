.PHONY:fmt vet build run

CURRENT_BRANCH := $(shell git branch --show-current)

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o build/http cmd/http/main.go

.PHONY: run
run:
	go run cmd/http/main.go

.PHONY: test
test:
	make test-resetdb && \
	STAGE=testing go test ./...

test-resetdb:
	STAGE=testing go run cmd/devcli/main.go dropdb && \
    STAGE=testing atlas migrate apply --url "postgres://dev:dev@devhost:5432/{{ cookiecutter.project_slug }}_test?sslmode=disable" --dir file://migrations

.PHONY: docs
docs:
	swag init -g cmd/http/main.go

dbreset:
	go run cmd/devcli/main.go dropdb && \
	atlas migrate apply --url "postgres://dev:dev@devhost:5432/{{ cookiecutter.project_slug }}?sslmode=disable" --dir file://migrations

dbdiff:
	atlas migrate diff --env gorm

dbapply:
	atlas migrate apply --url "postgres://dev:dev@devhost:5432/{{ cookiecutter.project_slug }}?sslmode=disable" --dir file://migrations

push:
	git checkout unstable \
	&& git pull origin unstable \
	&& git merge ${CURRENT_BRANCH} --no-edit \
	&& git push origin unstable \
	&& git checkout ${CURRENT_BRANCH}

worker:
	go run cmd/worker/main.go

mon:
	go run cmd/asynqmon/main.go