-include .make.env

DOCKER_REGISTRY ?= 
DOCKER_REPOSITORY ?=
VERSION ?= latest
SERVICE ?=
TYPE ?=

# ──────────────────────────────
# 🧪 Unit Tests
# ──────────────────────────────
test-unit:
	go test .../test/unit/...

# ──────────────────────────────
# Docker
# ──────────────────────────────
docker-build:
	@echo "🔨 Building $(SERVICE) $(if $(TYPE),($(TYPE)),(default)) image..."
	@DOCKERFILE=$(SERVICE)/Dockerfile$(if $(TYPE),.$(TYPE)) && \
	CONTEXT=$(SERVICE) && \
	IMAGE_NAME=$(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY)/$(SERVICE)$(if $(TYPE),-$(TYPE)):${VERSION} && \
	echo "🚀 Building image: $$IMAGE_NAME using $$DOCKERFILE" && \
	docker build -f $$DOCKERFILE -t $$IMAGE_NAME $$CONTEXT

docker-push:
	@echo "📦 Pushing $(SERVICE) $(if $(TYPE),($(TYPE)),(default)) image..."
	@IMAGE_NAME=$(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY)/$(SERVICE)$(if $(TYPE),-$(TYPE)):${VERSION} && \
	echo "🚀 Pushing image: $$IMAGE_NAME" && \
	docker push $$IMAGE_NAME
