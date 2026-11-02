.PHONY: check test docker-up docker-down

GOLANGCI_LINT_VERSION ?= v2.15.0

check:
ifeq ($(OS),Windows_NT)
	go test ./...
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	@GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION) ./scripts/install-golangci-lint.sh
	COVER_THRESHOLD=60.0 GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION) ./lint-project.sh
endif

test:
	go test ./...

docker-up:
	mkdir -p data
	chmod 777 data
	docker compose up -d
	./scripts/health-containers.sh

docker-down:
	docker compose down
