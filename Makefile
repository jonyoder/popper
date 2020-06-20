# note: call scripts from /scripts

# Tips:
# go mod tidy
# go mod vendor

BUILD_RUNNER=GOBIN=$(CURDIR)/bin

build:
	$(BUILD_RUNNER) go install -v ./...

test:
	go test -count=1 -v ./...

.PHONY: build test
