-include .make.env

# ──────────────────────────────
# 📦 Application Build Variables
# ──────────────────────────────
MODULE_NAME        ?= mandacode.com/accounts/auth
REGISTRY           ?=
APP_NAME           ?=
TAG                ?=
APP_IMAGE          := $(REGISTRY)/$(APP_NAME)

# ──────────────────────────────
# 🛠 Migration Image Variables
# ──────────────────────────────
MIGRATE_NAME       ?= ${APP_NAME}-migrate
MIGRATE_DEV_URL    ?=
MIGRATE_REGISTRY   ?= $(REGISTRY)
MIGRATE_TAG        ?= $(TAG)
MIGRATE_IMAGE      := $(MIGRATE_REGISTRY)/$(MIGRATE_NAME)

# ──────────────────────────────
# 📁 Directory Paths
# ──────────────────────────────
PROTO_SRC          := proto
DOCKER_CONTEXT     := ./docker
MOCK_DST_DIR       := test/mock
MOCK_TARGET_DIR    ?=

# ──────────────────────────────
# 🎯 Default Target
# ──────────────────────────────
all: build-app

# ──────────────────────────────
# 🚀 Run Application
# ──────────────────────────────
run:
	go run ./cmd/server

# ──────────────────────────────
# 🐳 Build & Push App Image
# ──────────────────────────────
build-app: check-tag
	docker build -t $(APP_IMAGE):$(TAG) -f $(DOCKER_CONTEXT)/app.Dockerfile .

push-app: build-app
	docker push $(APP_IMAGE):$(TAG)

# ──────────────────────────────
# 🔧 Migrations (Atlas)
# ──────────────────────────────
generate-migrate: check-migrate-name
	atlas migrate diff $(MIGRATE_NAME) \
		--dev-url "docker://postgres/17" \
		--dir "file://ent/migrate/migrations" \
		--to "ent://ent/schema"

apply-migrate:
	atlas migrate apply \
		--url $(MIGRATE_DEV_URL) \
		--dir "file://ent/migrate/migrations"

build-migrate: check-tag
	docker build -t $(MIGRATE_IMAGE):$(MIGRATE_TAG) -f $(DOCKER_CONTEXT)/migrate.Dockerfile .

push-migrate: build-migrate
	docker push $(MIGRATE_IMAGE):$(MIGRATE_TAG)

# ──────────────────────────────
# 🧬 Ent Schema Generation
# ──────────────────────────────
generate-ent:
	@echo "Generating Ent schema..."
	GOWORK=off go generate ./ent

# ──────────────────────────────
# 🧪 Unit Tests
# ──────────────────────────────
test-unit:
	go test ./test/unit/...

# ──────────────────────────────
# 🔍 Interface Mocks
# ──────────────────────────────
generate-mocks: check-mock-target-dir
	./scripts/generate-mocks.sh \
		-d $(MOCK_DST_DIR) \
		-t $(MOCK_TARGET_DIR) \
		-m $(MODULE_NAME)

# ──────────────────────────────
# 🛑 Validation Helpers
# ──────────────────────────────
check-tag:
ifndef TAG
	$(error TAG is required. Usage: make TAG=v1.2.3)
endif

check-migrate-name:
ifndef MIGRATE_NAME
	$(error MIGRATE_NAME is required. Usage: make MIGRATE_NAME=migrate)
endif

check-mock-target-dir:
ifndef MOCK_TARGET_DIR
	$(error MOCK_TARGET_DIR is required. Usage: make MOCK_TARGET_DIR=path/to/target)
endif

# ──────────────────────────────
# 🧾 PHONY Targets
# ──────────────────────────────
.PHONY: \
	all run \
	build-app push-app \
	generate-migrate apply-migrate build-migrate push-migrate \
	generate-ent test-unit \
	generate-mocks check-tag check-migrate-name check-mock-target-dir
