DOCKER_IMAGE     := jieht9u/prometheus-telegram-bot
DOCKERFILE_LOCAL := dockerfile
TEST_PKG         := msg 
TEST_DIR		 := ./vendor/src/app/#./vendor/src/app/...

VERSION          := $(shell git describe --tags --always --dirty="-dev")
DATE             := $(shell date -u '+%Y-%m-%d-%H:%M UTC')
VERSION_FLAGS    := -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'

bin_dir			 := $(CURDIR)/bin
dist_dir 		 := $(CURDIR)/dist

GITVERSION		 := 1.0.1

.PHONY: all 
all: test build

.PHONY: test
test: unit-test gometalinter

.PHONY: unit-test
unit-test:
	@echo "\n"
	@echo "[ Start Unit Testing ]"
	@for pkg in $(TEST_PKG) ; do \
		echo  "Run test [ $$pkg ] package" ; \
		GODEBUG=cgocheck=2 go test -race -v $(TEST_DIR)$$pkg; \
	done	

.PHONY: gometalinter
gometalinter: 
	@echo "\n"
	@echo "[ Start Gometalinter ] "
	@for pkg in $(TEST_PKG) ; do \
		echo  "Run test [ $$pkg ] package" ; \
		gometalinter $(TEST_DIR)$$pkg; \
	done	


.PHONY: build
build:
	@echo "\n"
	@echo "[ Building ]"
	$Q gb build


.PHONY: docker
docker: docker-build docker-push
	@echo "Success Docker"

.PHONY: docker-build
docker-build:
	@echo "Docker Build..."
	$Q docker build -t $(DOCKER_IMAGE):$(GITVERSION) --file=$(DOCKERFILE_LOCAL) .

.PHONY: docker-push
docker-push:
	@echo "Docker Push..."
	$Q docker tag $(DOCKER_IMAGE):$(GITVERSION) $(DOCKER_IMAGE):$(GITVERSION) 
	$Q docker push $(DOCKER_IMAGE):$(GITVERSION)


#TAG RELEASE
.PHONY: tag-release
tag-release: tag 
	$Q git push origin v$(GITVERSION)

# RELEASE
.PHONY: release
release: clean-dist tag 
	$Q goreleaser	

.PHONY: clean-dist
clean-dist:
	@echo "Removing distribution files"
	rm -rf $(dist_dir)

.PHONY: tags echo
tags:
	@echo "Listing tags..."
	$Q @git tag

echo:
	@echo "MESSAGE " $(MESSAGE)


.PHONY: tag
tag:
	@echo "Creating tag" $(GITVERSION)
	$Q @git tag -a v$(GITVERSION) -m $(GITVERSION)

.PHONY: td
td:
	@echo "Delete tag" $(GITVERSION)
	$Q @git tag -d v$(GITVERSION) 

#MERGE
.PHONY: merge
merge:
	$Q git checkout master
	$Q git merge develop
