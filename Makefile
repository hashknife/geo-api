GO ?= go
COVERAGEDIR = coverage
ifdef CIRCLE_ARTIFACTS
	COVERAGEDIR=$(CIRCLE_ARTIFACTS)/coverage
endif

.PHONY: test

TEST_LIST = $(foreach pkg, $(ALL_PACKAGES), $(pkg)_test)
COVER_LIST = $(foreach pkg, $(ALL_PACKAGES), $(pkg)_cover)
LDFLAGS = -ldflags "-X main.gitSHA=$(shell git rev-parse HEAD)"

deps:
	dep ensure

build:
	if [ ! -d bin ]; then mkdir bin; fi
	$(GO) build -v -o bin/geo-api $(LDFLAGS)

gen-mocks:
	mockery -dir=./endpoints/ -all
	mockery -dir=github.com/hashknife/common/services/ -all

test:
	if [ ! -d $(COVERAGEDIR) ]; then mkdir $(COVERAGEDIR); fi
	AWS_ACCESS_KEY_ID=1 AWS_SECRET_ACCESS_KEY=1 go test -v ./$* -race -cover -covermode=atomic -coverprofile=$(COVERAGEDIR)/$(subst /,_,$*).coverprofile
	$(GO) test -ldflags -s -v ./bindings -race -cover -coverprofile=$(COVERAGEDIR)/bindings.coverprofile
	$(GO) test -ldflags -s -v ./config -race -cover -coverprofile=$(COVERAGEDIR)/config.coverprofile
	$(GO) test -ldflags -s -v ./endpoints -race -cover -coverprofile=$(COVERAGEDIR)/endpoints.coverprofile

cover:
	go tool cover -html=$(COVERAGEDIR)/bindings.coverprofile -o $(COVERAGEDIR)/bindings.html
	go tool cover -html=$(COVERAGEDIR)/config.coverprofile -o $(COVERAGEDIR)/config.html
	go tool cover -html=$(COVERAGEDIR)/endpoints.coverprofile -o $(COVERAGEDIR)/endpoints.html

docs:
	@godoc -http=:6060 2>/dev/null &
	@printf "To view geo-api docs, point your browser to:\n"
	@printf "\n\thttp://127.0.0.1:6060/pkg/github.com/hashknife/geo-api/$(pkg)\n\n"
	@sleep 1
	@open "http://127.0.0.1:6060/pkg/github.com/hashknife/geo-api/$(pkg)"

tc: test cover

coveralls:
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=circle-ci -repotoken=$(COVERALLS_TOKEN); echo "Coveralls finished"

bench:
	go test -bench ./...

clean:
	$(GO) clean
	rm -f bin/geo-api
	rm -rf coverage/
