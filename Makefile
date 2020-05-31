GOCMD := go
GOBUILD := $(GOCMD) build
BUILD_DIR := $(PWD)
BIN_PATH := $(BUILD_DIR)/cmd
BINARY_NAME := ipfsmon
RELEASE_SUPPORT := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))/.release-support
MAKE_ALL := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))/.makeall

.release:
	@echo "release=0.0.0" >> .release
	@echo "tag=v0.0.0" >> .release
	@echo INFO: .release created
	@cat .release

tag-patch-release: VERSION := $(shell . $(RELEASE_SUPPORT); nextPatchLevel)
tag-patch-release: .release tag

tag-minor-release: VERSION := $(shell . $(RELEASE_SUPPORT); nextMinorLevel)
tag-minor-release: .release tag

tag-major-release: VERSION := $(shell . $(RELEASE_SUPPORT); nextMajorLevel)
tag-major-release: .release tag

patch-release: tag-patch-release check-release
	@echo $(VERSION)

minor-release: tag-minor-release check-release
	@echo $(VERSION)

major-release: tag-major-release check-release
	@echo $(VERSION)

tag: TAG=$(shell . $(RELEASE_SUPPORT); getTag $(VERSION))
tag: check-status
	@. $(RELEASE_SUPPORT) ; ! tagExists $(TAG) || (echo "ERROR: tag $(TAG) for version $(VERSION) already tagged in git" >&2 && exit 1) ;
	@. $(RELEASE_SUPPORT) ; setRelease $(VERSION)
	git add .
	git commit -m "bumped to version $(VERSION)" ;
	git tag $(TAG) ;
	@ if [ -n "$(shell git remote -v)" ] ; then git push --tags ; else echo 'no remote to push tags to' ; fi

check-status:
	@. $(RELEASE_SUPPORT) ; ! hasChanges || (echo "ERROR: there are still outstanding changes" >&2 && exit 1) ;

check-release: .release
	@. $(RELEASE_SUPPORT) ; tagExists $(TAG) || (echo "ERROR: version not yet tagged in git. make [minor,major,patch]-release." >&2 && exit 1) ;
	@. $(RELEASE_SUPPORT) ; ! differsFromRelease $(TAG) || (echo "ERROR: current directory differs from tagged $(TAG). make [minor,major,patch]-release." ; exit 1)


showver: .release
	@. $(RELEASE_SUPPORT); getTag

build: VERSION := $(shell . $(RELEASE_SUPPORT); getVersion)
build:
	@-echo "Building binary..."
	@-$(GOBUILD) -ldflags="-X 'main.CurrentVersion=$(VERSION)'" -o $(BIN_PATH)/$(BINARY_NAME)

install: VERSION := $(shell . $(RELEASE_SUPPORT); getVersion)
install:
	@-echo "Installing binary..."
	@-$(GOBUILD) -ldflags="-X 'main.CurrentVersion=$(VERSION)'" -o $(GOBIN)/$(BINARY_NAME)

all: VERSION := $(shell . $(RELEASE_SUPPORT); getVersion)
all:
	@-echo "Building binaries for all platforms"
	@. $(MAKE_ALL) ; BIN_PATH=$(BIN_PATH) BINARY_NAME=$(BINARY_NAME) makeall $(VERSION)
