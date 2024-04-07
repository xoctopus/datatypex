tidy:
	go mod tidy
cover: tidy
	go test -v -race -failfast -parallel 1 -gcflags="all=-N -l" ./... -covermode=atomic -coverprofile cover.out
test: tidy
	go test -v -race -failfast -parallel 1 -gcflags="all=-N -l" ./...
