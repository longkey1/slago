.DEFAULT_GOAL := help

export GO_VERSION=$(shell grep "^go " go.mod | sed 's/^go //')
export PRODUCT_NAME := slago

.PHONY: build
build: ## Build the application to ./bin/
	go build -o ./bin/$(PRODUCT_NAME) ./cmd/slago

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf ./bin .cache

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: tools
tools: ## Install tools
	go install github.com/goreleaser/goreleaser/v2@latest

# === Version management ===
VERSION := $(shell v=$$(git tag --sort=-v:refname | head -n1 2>/dev/null); [ -n "$$v" ] && echo "$$v" || echo "v0.0.0")
MAJOR := $(shell echo $(VERSION) | sed 's/v//' | cut -d. -f1)
MINOR := $(shell echo $(VERSION) | sed 's/v//' | cut -d. -f2)
PATCH := $(shell echo $(VERSION) | sed 's/v//' | cut -d. -f3)

.PHONY: release
release: ## Release a new version (usage: make release type=patch|minor|major [dryrun=false])
	@type=$${type:-patch}; \
	dryrun=$${dryrun:-true}; \
	case "$$type" in \
		major) new_version="v$$(($(MAJOR)+1)).0.0" ;; \
		minor) new_version="v$(MAJOR).$$(($(MINOR)+1)).0" ;; \
		patch) new_version="v$(MAJOR).$(MINOR).$$(($(PATCH)+1))" ;; \
		*) echo "Invalid type: $$type. Use patch, minor, or major."; exit 1 ;; \
	esac; \
	echo "Current version: $(VERSION)"; \
	echo "New version: $$new_version"; \
	if [ "$$dryrun" = "false" ]; then \
		git tag -a "$$new_version" -m "Release $$new_version"; \
		git push origin "$$new_version"; \
		echo "Released $$new_version"; \
	else \
		echo "[dryrun] Would create and push tag: $$new_version"; \
	fi

.PHONY: re-release
re-release: ## Re-release an existing version (usage: make re-release [tag=<tag>] [dryrun=false])
	@tag=$${tag:-$(VERSION)}; \
	dryrun=$${dryrun:-true}; \
	echo "Re-releasing: $$tag"; \
	if [ "$$dryrun" = "false" ]; then \
		gh release delete "$$tag" --yes 2>/dev/null || true; \
		git push origin --delete "$$tag" 2>/dev/null || true; \
		git tag -d "$$tag" 2>/dev/null || true; \
		git tag -a "$$tag" -m "Release $$tag"; \
		git push origin "$$tag"; \
		echo "Re-released $$tag"; \
	else \
		echo "[dryrun] Would delete and recreate tag: $$tag"; \
	fi

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
