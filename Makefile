tidy:
	go mod tidy
cover: tidy
	go test -race -failfast -parallel 1 -gcflags="all=-N -l" ./... -covermode=atomic -coverprofile cover.out
test: tidy
	go test -race -failfast -parallel 1 -gcflags="all=-N -l" ./...

report:
	@echo ">>>static checking"
	@go vet ./...
	@echo "done\n"
	@echo ">>>detecting ineffectual assignments"
	@ineffassign ./...
	@echo "done\n"
	@echo ">>>detecting icyclomatic complexities over 10 and average"
	@gocyclo -over 10 -avg -ignore '_test|vendor' .
	@echo "done\n"

