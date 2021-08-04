PACKAGES := $(shell go list ./... | grep -v '/lib/' | grep -v '/vendor')

.PHONY: test
test:
	go test -coverprofile=profile.cov -covermode=atomic -coverpkg=github.com/bjartek/go-with-the-flow/v2/gwtf -v ./...

.PHONY: cover
cover: test
	 go tool cover -html=profile.cov

golint:
	@echo "Running golint"
	@golint --set_exit_status $(PACKAGES)