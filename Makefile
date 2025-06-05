PACKAGES=$(shell go list ./... | grep -E -v 'pb$|testdata|mock|proto|example')
MOD=$(shell cat go.mod | grep ^module -m 1 | awk '{ print $$2; }' || '')
MOD_NAME=$(shell basename $(MOD))

GOTEST=xgo
GOBUILD=go
GOFMT=goimports-reviser

# dependencies
GOIMPORTS=$(shell type goimports-reviser > /dev/null 2>&1 && echo $$?)
XGO=$(shell type xgo > /dev/null 2>&1 && echo $$?)

# git info
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --tags --abbrev=0)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_AT=$(shell date "+%Y%m%d%H%M%S")

show:
	@echo "packages:"
	@for item in $(PACKAGES); do echo "\t$$item"; done
	@echo "module:"
	@echo "\t$(MOD)"
	@echo "tools:"
	@echo "\tbuild=$(GOBUILD) test=$(GOTEST) fmt=$(GOFMT) xgo=$(XGO)"
	@echo "git:"
	@echo "\tcommit_id=$(GIT_COMMIT)\n\ttag=$(GIT_TAG)\n\tbranch=$(GIT_BRANCH)\n\tbuild_time=$(BUILD_AT)\n\tname=$(MOD_NAME)"

# install dependencies
dep:
	@if [ "${GOIMPORTS}" != "0" ]; then \
		echo "installing goimports-reviser for format sources"; \
		go install -v github.com/incu6us/goimports-reviser/v3@latest; \
	fi
	@if [ "${GOTEST}" = "xgo" ] && [ "${XGO}" != "0" ]; then \
		echo "installing xgo for unit test"; \
		go install github.com/xhd2015/xgo/cmd/xgo@latest; \
	fi

update:
	@echo "==> installing goimports-reviser for format sources"
	@go install -v github.com/incu6us/goimports-reviser/v3@latest
	@if [ "${GOTEST}" == "xgo" ]; then \
		echo "==> installing xgo for unit test"; \
		go install github.com/xhd2015/xgo/cmd/xgo@latest; \
	fi
	@echo "==> updating go module dependencies"
	@go get -u all

tidy:
	@echo "==> tidy"
	@go mod tidy

cover: show dep tidy
	@echo "==> run unit test with coverage"
	@$(GOTEST) test -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES} -covermode=count -coverprofile=cover.out

view_cover: dep tidy
	@echo "==> run unit test with coverage and view"
	@$(GOTEST) test -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES} -covermode=count -coverprofile=cover.out && go tool cover -html cover.out

test: dep tidy
	@echo "==> run unit test"
	@$(GOTEST) test -race -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES}

benchmark: dep tidy
	@echo "==> run benchmark"
	@$(GOBUILD) test -bench=. -benchmem ${PACKAGES}

vet:
	@go vet ${PACKAGES}

fmt: clean
	@echo "==> format code"
	@if [ "${GOFMT}" == "goimports-reviser" ]; then \
		goimports-reviser -rm-unused -set-alias -output write -imports-order 'std,general,company,project' -excludes '.git/,.xgo/,*.pb.go,*_generated.go' ./...; \
	else \
		go fmt ./...; \
	fi

report:
	@echo ">>>static checking"
	@go vet ./...
	@echo "done\n"
	@echo ">>>detecting ineffectual assignments"
	@ineffassign ./...
	@echo "done\n"
	@echo ">>>detecting icyclomatic complexities over 10 and average"
	@gocyclo -over 10 -avg -ignore '_test|vendor' . || true
	@echo "done\n"

clean:
	@find . -name cover.out | xargs rm -rf
	@find . -name .xgo | xargs rm -rf
	@rm -rf build/*

check: show tidy update vet cover fmt clean

benchmark_snowflake:
	@cd snowflake/internal
	@go test -bench=BenchmarkSnowflake_ID -benchtime=30000x -unit=1 ## modify unit to assign factory
