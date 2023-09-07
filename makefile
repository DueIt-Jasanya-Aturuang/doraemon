.PHONY: test

test:
	@echo "\n---- RUNNING UNIT TEST TESTS/UNIT FOLDER ----\n"
	go test -v ./tests/unit/
	go clean -testcache

linter:
	@echo "---- START LINTER GOLANG ----\n"
	golangci-lint version
	golangci-lint run -c .golangci.yml ./...

go-generate:
	go generate ./...