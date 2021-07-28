.PHONY: test
test:
	go test -coverprofile=profile.cov -covermode=atomic -coverpkg=github.com/bjartek/go-with-the-flow/v2/gwtf ./...

.PHONY: cover
cover: test
	 go tool cover -html=profile.cov

