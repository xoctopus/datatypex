PACKAGES=$(shell go list ./... | grep -E -v 'pb$|testdata|mock|proto|example')
MOD=$(shell cat go.mod | grep ^module -m 1 | awk '{ print $$2; }' || '')
MOD_NAME=$(shell basename $(MOD))

GOTEST=xgo
GOBUILD=go

# dependencies
DEP_FMT=$(shell type goimports-reviser > /dev/null 2>&1 && echo $$?)
DEP_XGO=$(shell type xgo > /dev/null 2>&1 && echo $$?)

# git info
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --tags --abbrev=0)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_AT=$(shell date "+%Y%m%d%H%M%S")

show:
	@echo "packages:"
	@for item in $(PACKAGES); do echo "\t$$item"; done
	@echo "module:"
	@echo "\tpath=$(MOD)"
	@echo "\tmodule=$(MOD_NAME)"
	@echo "tools:"
	@echo "\tbuild=$(GOBUILD)"
	@echo "\ttest=$(GOTEST)"
	@echo "\tgoimports-reviser=$(shell which goimports-reviser)"
	@echo "\txgo=$(shell which xgo)"
	@echo "git info:"
	@echo "\tcommit_id=$(GIT_COMMIT)\n\ttag=$(GIT_TAG)\n\tbranch=$(GIT_BRANCH)\n\tbuild_time=$(BUILD_AT)\n\tname=$(MOD_NAME)"

# install dependencies
dep:
	@echo "==> installing dependencies"
	@if [ "${DEP_FMT}" != "0" ]; then \
		echo "\tgoimports-reviser for format sources"; \
		go install github.com/incu6us/goimports-reviser/v3@latest; \
	fi
	@if [ "${GOTEST}" = "xgo" ] && [ "${DEP_XGO}" != "0" ]; then \
		echo "\txgo for unit test"; \
		go install github.com/xhd2015/xgo/cmd/xgo@latest; \
	fi

upgrade-dep:
	@echo "==> upgrading dependencies"
	@echo "\tupgrading goimports-reviser for format sources"
	@go install github.com/incu6us/goimports-reviser/v3@latest
	@if [ "${GOTEST}" = "xgo" ]; then \
		echo "\tupgrading xgo for unit test"; \
		go install github.com/xhd2015/xgo/cmd/xgo@latest; \
	fi

update:
	@go get -u all

tidy:
	@echo "==> tidy"
	@go mod tidy

test: dep tidy
	@echo "==> run unit test"
	@if [ "${GOTEST}" = "xgo" ]; then \
		$(GOTEST) test -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES}; # xgo mock functions may cause data race \
	else \
		$(GOTEST) test -race -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES}; \
	fi

cover: dep tidy
	@echo "==> run unit test with coverage"
	@$(GOTEST) test -failfast -parallel 1 -gcflags="all=-N -l" ${PACKAGES} -covermode=count -coverprofile=cover.out

view_cover: cover
	@echo "==> run unit test with coverage and view"
	@go tool cover -html cover.out

vet:
	@go vet ${PACKAGES}

fmt: dep clean
	@echo "==> format code"
	@goimports-reviser -rm-unused -set-alias -output write -imports-order 'std,general,company,project' -excludes '.git/,.xgo/,*.pb.go,*_generated.go' ./...

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
