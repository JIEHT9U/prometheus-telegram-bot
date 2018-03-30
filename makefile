DOCKER_IMAGE     := jieht9u/prometheus-telegram-bot
bin_dir			 := $(CURDIR)/bin
VER				 := 1.0.0
TAG              := v$(VERSION)
DOCKERFILE_LOCAL := dockerfile
TEST_PKG         := msg 
TEST_DIR		 := ./vendor/src/app/#./vendor/src/app/...


VERSION          := $(shell git describe --tags --always --dirty="-dev")
DATE             := $(shell date -u '+%Y-%m-%d-%H:%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'


.PHONY: all 
all: test build

.PHONY: test
test: unit-test gometalinter

.PHONY: unit-test
unit-test:
	@echo "Unit Testing..."
	@for pkg in $(TEST_PKG) ; do \
		GODEBUG=cgocheck=2 go test -race -v $(TEST_DIR)$$pkg; \
	done	

.PHONY: gometalinter
gometalinter: gometalinter
	@echo "Gometalinter Run..."
	@for pkg in $(TEST_PKG) ; do \
		gometalinter $(TEST_DIR)$$pkg; \
	done	


.PHONY: build
build:
	@echo "Building..."
	$Q gb build


.PHONY: docker
docker: docker-build docker-push
	@echo "Success Docker"

.PHONY: docker-build
docker-build:
	@echo "Docker Build..."
	$Q docker build -t $(DOCKER_IMAGE):$(VER) --file=$(DOCKERFILE_LOCAL) .

.PHONY: docker-push
docker-push:
	@echo "Docker Push..."
	$Q docker tag $(DOCKER_IMAGE):$(VER) $(DOCKER_IMAGE):$(VER) 
	$Q docker push $(DOCKER_IMAGE):$(VER)


