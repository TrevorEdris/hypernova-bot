SHELL ?= /bin/bash
export REGISTRY ?= ${DOCKER_REGISTRY}
export IMAGEORG ?= tedris
export IMAGEBASE ?= hypernova
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
	docker-compose -f ${DEV_DOCKER_COMPOSE} up --build -d
	make -s dev-logs

.PHONY: dev-down
dev-down: ## Stop all containers
	docker-compose -f ${DEV_DOCKER_COMPOSE} down

.PHONY: dev-restart
dev-restart: ## Restart all containers
	docker-compose -f ${DEV_DOCKER_COMPOSE} restart

.PHONY: dev-logs
dev-logs: ## Print logs in stdout
	docker-compose -f ${DEV_DOCKER_COMPOSE} logs -f api bot populate

# -----------------------------[ Build ]-----------------------------

.PHONY: build
build: submodules version build-api build-bot ## Build and tag the docker container for the API and Bot

.PHONY: build-api
build-api:
	@docker build -f deployments/container/api.Dockerfile -t ${IMAGEORG}/${IMAGEBASE}-api:${VERSION} --target builder .
	@docker tag ${IMAGEORG}/${IMAGEBASE}-api:${VERSION} ${IMAGEORG}/${IMAGEBASE}-api:latest
	@docker tag ${IMAGEORG}/${IMAGEBASE}-api:${VERSION} ${IMAGEORG}/${IMAGEBASE}-api-build:latest

.PHONY: build-bot
build-bot:
	@docker build -f deployments/container/bot.Dockerfile -t ${IMAGEORG}/${IMAGEBASE}-bot:${VERSION} --target builder .
	@docker tag ${IMAGEORG}/${IMAGEBASE}-bot:${VERSION} ${IMAGEORG}/${IMAGEBASE}-bot:latest
	@docker tag ${IMAGEORG}/${IMAGEBASE}-bot:${VERSION} ${IMAGEORG}/${IMAGEBASE}-bot-build:latest

.PHONY: build-populate
build-populate: ## Build and tag the container that populates data for the local dev environment
	@docker build -f deployments/container/populate.Dockerfile -t ${IMAGEORG}/${IMAGEBASE}-populate:latest .

# -----------------------------[ Test ]------------------------------

.PHONY: test
test: build test-unit ## Run full test suite

.PHONY: test-unit
test-unit: ## Run unit tests
	@test/test_unit

# -----------------------------[ Publish ]---------------------------

.PHONY: finalize
finalize: test finalize-api finalize-bot ## Build, test, and tag the docker container with the finalized tag (typically, the full docker registery will be tagged here)

.PHONY: finalize-api
finalize-api:
	@docker build -f container/Dockerfile -t ${IMAGEORG}/${IMAGEBASE}-api:${VERSION} .
	@docker tag ${IMAGEORG}/${IMAGEBASE}-api:${VERSION} ${IMAGEORG}/${IMAGEBASE}-api:latest

.PHONY: finalize-bot
finalize-bot:
	@docker build -f container/Dockerfile -t ${IMAGEORG}/${IMAGEBASE}-bot:${VERSION} .
	@docker tag ${IMAGEORG}/${IMAGEBASE}-bot:${VERSION} ${IMAGEORG}/${IMAGEBASE}-bot:latest

.PHONY: publish-only
publish-only: publish-only-api publish-only-bot ## Push the tagged docker image to the docker registry

.PHONY: publish-only-api
publish-only-api: ## Push the tagged docker image to the docker registry
	@docker tag ${IMAGEORG}/${IMAGEBASE}-api:${VERSION} ${REGISTRY}${IMAGEORG}/${IMAGEBASE}-api:${VERSION}
	@docker push ${REGISTRY}${IMAGEORG}/${IMAGEBASE}-api:${VERSION}

.PHONY: publish-only-bot
publish-only-bot: ## Push the tagged docker image to the docker registry
	@docker tag ${IMAGEORG}/${IMAGEBASE}-bot:${VERSION} ${REGISTRY}${IMAGEORG}/${IMAGEBASE}-bot:${VERSION}
	@docker push ${REGISTRY}${IMAGEORG}/${IMAGEBASE}-bot:${VERSION}

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
copy-binary: build copy-binary-api copy-binary-bot ## Create a temporary container based on the "-build" image and copy the binary out of the container

.PHONY: copy-binary-api
copy-binary-api: ## Create a temporary container based on the "-build" image and copy the binary out of the container
	@docker create --name ${IMAGEBASE}-api-${GIT_HASH} ${IMAGEORG}/${IMAGEBASE}-api-build:${VERSION}
	@docker cp ${IMAGEBASE}-api-${GIT_HASH}:/src/the-binary ./the-binary
	@docker rm ${IMAGEBASE}-api-${GIT_HASH}

.PHONY: copy-binary-bot
copy-binary-bot: ## Create a temporary container based on the "-build" image and copy the binary out of the container
	@docker create --name ${IMAGEBASE}-bot-${GIT_HASH} ${IMAGEORG}/${IMAGEBASE}-bot-build:${VERSION}
	@docker cp ${IMAGEBASE}-bot-${GIT_HASH}:/src/the-binary ./the-binary
	@docker rm ${IMAGEBASE}-bot-${GIT_HASH}
