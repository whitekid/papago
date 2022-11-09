TEST_FLAGS ?= -v

.PHONY: test dep tidy

test:
	@go test ${TEST_FLAGS} ./...

# update modules & tidy
dep:
	@rm -f go.mod go.sum
	@go mod init papago

	@$(MAKE) tidy

tidy:
	@go mod tidy -v
