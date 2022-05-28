SHELL ?= /bin/bash
export REGISTRY ?= ${DOCKER_REGISTRY}
export IMAGEORG ?= tedris
export IMAGE ?= api-template
export VERSION ?= $(shell printf "`./tools/version`${VERSION_SUFFIX}")
export GIT_HASH =$(shell git rev-parse --short HEAD)
export DEV_DOCKER_COMPOSE ?= deployments/local/docker-compose.dev.yaml

# Blackbox files that need to be decrypted.
clear_files=$(shell blackbox_list_files)
encrypt_files=$(patsubst %,%.gpg,${clear_files})

# =========================[ Common Targets ]========================
# These are targets that almost certainly will not need to be changed
# as they are common to nearly all repos.
# ===================================================================

.PHONY: all
all: build

.PHONY: help
help: ## List of available commands
	@echo "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\033[36m\1\\033[m:\2/' | column -c2 -t -s :)"

# -------------------------[ General Tools ]-------------------------

.PHONY: clear
clear: ${clear_files}

${clear_files}: ${encrypt_files}
	@blackbox_decrypt_all_files

.PHONY: decrypt
decrypt: ${clear_files} ## Decrypt all .gpg files registered in .blackbox/blackbox-files.txt

.PHONY: encrypt
encrypt: ${encrypt_files} ## Encrypt all files registered in .blackbox/blackbox-files.txt
	blackbox_edit_end $^

.PHONY: submodules
submodules: ## Recursively init all submodules in the repo
	@git submodule update --init --recursive || printf "\nWarning: Could not pull submodules\n"

.PHONY: version
version: submodules tools/version ## Automatically calculate the version
	@echo ${VERSION}

# =========================[ Custom Targets ]========================
# These are targets that _may_ need to be customized to the specific
# project implemented in the repo.
# ===================================================================

# ---------------------------[ Local App ]---------------------------
.PHONY: dev
dev: build-populate ## Run the API locally and print logs to stdout
	docker-compose -f ${DEV_DOCKER_COMPOSE} up -d
	make -s dev-logs

.PHONY: dev-down
dev-down: ## Stop all containers
	docker-compose -f ${DEV_DOCKER_COMPOSE} down

.PHONY: dev-restart
dev-restart: ## Restart all containers
	docker-compose -f ${DEV_DOCKER_COMPOSE} restart

.PHONY: dev-logs
dev-logs: ## Print logs in stdout
	docker-compose -f ${DEV_DOCKER_COMPOSE} logs -f app populate

# -----------------------------[ Build ]-----------------------------

.PHONY: build
build: submodules version ## Build and tag the docker container for the API
	@docker build -f deployments/container/Dockerfile -t ${IMAGEORG}/${IMAGE}:${VERSION} --target builder .
	@docker tag ${IMAGEORG}/${IMAGE}:${VERSION} ${IMAGEORG}/${IMAGE}:latest
	@docker tag ${IMAGEORG}/${IMAGE}:${VERSION} ${IMAGEORG}/${IMAGE}-build:latest

.PHONY: build-populate
build-populate: ## Build and tag the container that populates data for the local dev environment
	@docker build -f deployments/container/populate.Dockerfile -t ${IMAGEORG}/${IMAGE}-populate:latest .

# -----------------------------[ Test ]------------------------------

.PHONY: test
test: build test-unit ## Run full test suite

.PHONY: test-unit
test-unit: ## Run unit tests
	@test/test_unit

# -----------------------------[ Publish ]---------------------------

.PHONY: finalize
finalize: test ## Build, test, and tag the docker container with the finalized tag (typically, the full docker registery will be tagged here)
	@docker build -f container/Dockerfile -t ${IMAGEORG}/${IMAGE}:${VERSION} .
	@docker tag ${IMAGEORG}/${IMAGE}:${VERSION} ${IMAGEORG}/${IMAGE}:latest

.PHONY: publish-only
publish-only: ## Push the tagged docker image to the docker registry
	@docker tag ${IMAGEORG}/${IMAGE}:${VERSION} ${REGISTRY}${IMAGEORG}/${IMAGE}:${VERSION}
	@docker push ${REGISTRY}${IMAGEORG}/${IMAGE}:${VERSION}

.PHONY: publish
publish: finalize publish-only ## Finalize and publish the docker container

# -----------------------------[ Deploy ]----------------------------

.PHONY: kube-deploy-only
kube-deploy-only: decrypt ## Fill out the .yaml.tmpl files and apply them to the specified namespace
	@deployments/kube/deploy

.PHONY: kube-deploy
kube-deploy: publish kube-deploy-only ## Build, test, finalize, publish, and then deploy the docker container to kube

# ----------------------------[ Release ]----------------------------
# TODO

# -----------------------------[ Other ] ----------------------------

.PHONY: copy-binary
copy-binary: build ## Create a temporary container based on the "-build" image and copy the binary out of the container
	@docker create --name ${IMAGE}-${GIT_HASH} ${IMAGEORG}/${IMAGE}-build:${VERSION}
	@docker cp ${IMAGE}-${GIT_HASH}:/src/the-binary ./the-binary
	@docker rm ${IMAGE}-${GIT_HASH}
