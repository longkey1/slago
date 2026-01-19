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
	echo "Next version: $$new_version"; \
	if [ "$$dryrun" = "false" ]; then \
		echo "Creating new tag $$new_version..."; \
		git push origin master --no-verify --force-with-lease; \
		git tag -a "$$new_version" -m "Release of $$new_version"; \
		git push origin "$$new_version" --no-verify --force-with-lease; \
		echo "Tag $$new_version has been created and pushed"; \
		echo "GitHub Actions will build the release binary automatically"; \
	else \
		echo "[DRY RUN] Showing what would be done..."; \
		echo "Would push to origin/master"; \
		echo "Would create tag: $$new_version"; \
		echo "Would push tag to origin: $$new_version"; \
		echo ""; \
		echo "To execute this release, run:"; \
		echo "  make release type=$$type dryrun=false"; \
		echo "Dry run complete."; \
	fi

.PHONY: re-release
re-release: ## Re-release an existing version (usage: make re-release [tag=<tag>] [dryrun=false])
	@tag=$${tag:-$(VERSION)}; \
	dryrun=$${dryrun:-true}; \
	echo "Target tag: $$tag"; \
	if [ "$$dryrun" = "false" ]; then \
		echo "Deleting GitHub release..."; \
		gh release delete "$$tag" -y; \
		echo "Deleting local tag..."; \
		git tag -d "$$tag"; \
		echo "Deleting remote tag..."; \
		git push origin ":refs/tags/$$tag" --no-verify --force; \
		echo "Recreating tag on HEAD..."; \
		git tag -a "$$tag" -m "Release $$tag"; \
		echo "Pushing tag to origin..."; \
		git push origin "$$tag" --no-verify --force-with-lease; \
		echo "Recreating GitHub release..."; \
		gh release create "$$tag" --title "$$tag" --notes "Re-release of $$tag"; \
		echo "GitHub Actions will build the release binary automatically"; \
		echo "Done!"; \
	else \
		echo "[DRY RUN] Showing what would be done..."; \
		echo "Would delete release: $$tag"; \
		echo "Would delete local tag: $$tag"; \
		echo "Would delete remote tag: $$tag"; \
		echo "Would create new tag at HEAD: $$tag"; \
		echo "Would push tag to origin: $$tag"; \
		echo "Would create new release for: $$tag"; \
		echo ""; \
		echo "To execute this re-release, run:"; \
		echo "  make re-release tag=$$tag dryrun=false"; \
		echo "Dry run complete."; \
	fi

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
