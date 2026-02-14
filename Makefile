REACT_UI_DIR := react_ui
CMD_DIR := cmd
SAMPLE_DATA := sample_data/AAPL.parquet

# Default target
.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build-ui
build-ui: ## Build the React UI into internal/web/dist
	cd $(REACT_UI_DIR) && npm run build

.PHONY: build
build: build-ui ## Build everything (UI + Go binary)
	cd $(CMD_DIR) && go build -o ../bin/cibo .

.PHONY: run-tui
run-tui: build-ui ## Build UI then launch the TUI
	cd $(CMD_DIR) && go run .

.PHONY: run-tui-mock
run-tui-mock: build-ui ## Build UI then launch the TUI with mock API
	cd $(CMD_DIR) && go run . -mockAPI

.PHONY: run-web-sample
run-web-sample: build-ui ## Build UI then launch in standalone web mode (uses sample data)
	cd $(CMD_DIR) && go run . -webMode ../$(SAMPLE_DATA)

.PHONY: run-ui-dev
run-ui-dev: ## Start the Vite dev server with hot reload (no Go rebuild)
	cd $(REACT_UI_DIR) && npm run dev

.PHONY: test-go
test-go: ## Run all Go tests
	go test ./... -v

.PHONY: test-ui
test-ui: ## Run React unit tests
	cd $(REACT_UI_DIR) && npm test

.PHONY: test
test: test-go test-ui ## Run all tests (Go + React)

.PHONY: lint-go
lint-go: ## Run staticcheck on Go code
	staticcheck ./...

.PHONY: lint-ui
lint-ui: ## Run ESLint on React code
	cd $(REACT_UI_DIR) && npm run lint

.PHONY: lint
lint: lint-go lint-ui ## Run all linters

.PHONY: install-ui
install-ui: ## Install React dependencies
	cd $(REACT_UI_DIR) && npm install

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf internal/web/dist/*
	rm -rf bin/
	@# Restore the .gitkeep so go:embed doesn't break
	touch internal/web/dist/.gitkeep
