GO=go
GOBIN=${GOPATH}/bin

lint: go.lint.verify
	@echo "===========> Start running linters"
	@echo "=================> Running golangci-lint to check go files"
	@$(GOBIN)/golangci-lint run --fix
	@echo "=================> Running logcheck script to check log and err"
	@script/linter-logcheck.sh $(PWD)
	@echo "===========> Finish running linters"

test:
	@sh ${PWD}/script/gotest.sh

.PHONY: go.lint.verify
go.lint.verify:
ifeq (,$(wildcard $(GOBIN)/golangci-lint))
	@echo "===========> Installing golangci-lint"
	@GO111MODULE=on $(GO) get github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint
endif

