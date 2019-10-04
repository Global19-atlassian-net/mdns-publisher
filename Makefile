.DEFAULT_GOAL:=help

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./pkg/... ./cmd/...
	git diff --exit-code

.PHONY: vet
vet: ## Run go vet against code
	go vet ./pkg/... ./cmd/...
