lint-imports:
	@find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" | while read -r file; do \
		goimports-reviser -company-prefixes github.com/0xPellNetwork/pellapp-sdk -rm-unused -format "$$file"; \
	done

#? lint: Run latest golangci-lint linter
lint:
	@echo "--> Running linter"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run
.PHONY: lint

#? vulncheck: Run latest govulncheck
vulncheck:
	@echo "--> Running go vuln check"
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...
.PHONY: vulncheck

test:
	@echo "--> Running tests"
	@go test -v ./...
.PHONY: test


pre-commit:
	@make lint-imports
	@make lint
	@make vulncheck
	@make test

.PHONY: proto

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/test/service.proto
