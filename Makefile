# BUILD ===========================
current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)

# set dummy version
ifeq ($?, 0)
else
	git_description = "demo-pre"
endif

linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'

	cd backend && go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	# GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api

