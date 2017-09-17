GO ?= go
COVERAGEDIR = coverage
ifdef CIRCLE_ARTIFACTS
	COVERAGEDIR=$(CIRCLE_ARTIFACTS)/coverage
endif

.PHONY: test

ALL_PACKAGES = \
	logger \
	metrics \
	middleware \
	models \
	notifications \
	reports \
	services \
	utils

TEST_LIST = $(foreach pkg, $(ALL_PACKAGES), $(pkg)_test)
COVER_LIST = $(foreach pkg, $(ALL_PACKAGES), $(pkg)_cover)

deps:
	dep ensure

gen-mocks:
	mockery -dir=./logger/ -all
	mockery -dir=./metrics/ -all
	mockery -dir=./models/ -all
	mockery -dir=./middleware/ -all
	mockery -dir=./notifications/ -all
	mockery -dir=./reports/ -all
	mockery -dir=./services/ -all
	mockery -dir=./utils/ -all

.PHONY: test
test: $(TEST_LIST)

.PHONY: assert-no-diff
assert-no-diff:
	@./assert-no-diff.sh

$(TEST_LIST): %_test:
	@if [ ! -d coverage ]; then mkdir coverage; fi
	@AWS_ACCESS_KEY_ID=1 AWS_SECRET_ACCESS_KEY=1 go test -v ./$* -race -cover -coverprofile=$(COVERAGEDIR)/$(subst /,_,$*).coverprofile


.PHONY: cover
cover: $(COVER_LIST)

$(COVER_LIST): %_cover:
	$(GO) tool cover -html=$(COVERAGEDIR)/$(subst /,_,$*).coverprofile -o $(COVERAGEDIR)/$(subst /,_,$*).html

docs:
	@godoc -http=:6060 2>/dev/null &
	@printf "To view geo-api docs, point your browser to:\n"
	@printf "\n\thttp://127.0.0.1:6060/pkg/github.com/hashknife/common/$(pkg)\n\n"
	@sleep 1
	@open "http://127.0.0.1:6060/pkg/github.com/hashknife/common/$(pkg)"

tc: test cover

coveralls:
	gover $(COVERAGEDIR) $(COVERAGEDIR)/coveralls.coverprofile
	goveralls -coverprofile=$(COVERAGEDIR)/coveralls.coverprofile -service=circle-ci -repotoken=$(COVERALLS_TOKEN); echo "Coveralls finished"

bench:
	go test -bench ./...

clean:
	$(GO) clean
	rm -rf coverage/
