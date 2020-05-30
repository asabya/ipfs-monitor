GOCMD := go
GOBUILD := $(GOCMD) build
BUILD_DIR := $(PWD)
BIN_PATH := $(BUILD_DIR)/cmd
BINARY_NAME := mipfs
BUILD_SUPPORT := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))/.build-support
MAKE_ALL := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))/.makeall


.release:
	@echo "release=0.0.0" > .release
	@echo "tag=v0.0.0" >> .release
	@echo INFO: .release created
	@cat .release

tag-patch-release: VERSION := $(shell . $(BUILD_SUPPORT); nextPatchLevel)
tag-patch-release: .release tag

tag-minor-release: VERSION := $(shell . $(BUILD_SUPPORT); nextMinorLevel)
tag-minor-release: .release tag

tag-major-release: VERSION := $(shell . $(BUILD_SUPPORT); nextMajorLevel)
tag-major-release: .release tag

patch-release: tag-patch-release release
	@echo $(VERSION)

minor-release: tag-minor-release release
	@echo $(VERSION)

major-release: tag-major-release release
	@echo $(VERSION)

tag: TAG=$(shell . $(BUILD_SUPPORT); getTag $(VERSION))
tag: check-status check-release
	@. $(BUILD_SUPPORT) ; ! tagExists $(TAG) || (echo "ERROR: tag $(TAG) for version $(VERSION) already tagged in git" >&2 && exit 1) ;
	@. $(BUILD_SUPPORT) ; setRelease $(VERSION)
	git add .
	git commit -m "bumped to version $(VERSION)" ;
	git tag $(TAG) ;
#	@ if [ -n "$(shell git remote -v)" ] ; then git push --tags ; else echo 'no remote to push tags to' ; fi

check-status:
	@. $(BUILD_SUPPORT) ; ! hasChanges || (echo "ERROR: there are still outstanding changes" >&2 && exit 1) ;

check-release: .release
	@. $(BUILD_SUPPORT) ; tagExists $(TAG) || (echo "ERROR: version not yet tagged in git. make [minor,major,patch]-release." >&2 && exit 1) ;
	@. $(BUILD_SUPPORT) ; ! differsFromRelease $(TAG) || (echo "ERROR: current directory differs from tagged $(TAG). make [minor,major,patch]-release." ; exit 1)


showver: .release
	@. $(BUILD_SUPPORT); getTag

build:
	@-echo "Building binary..."
	@-$(GOBUILD) -o $(BIN_PATH)/$(BINARY_NAME)

install:
	@-echo "Installing binary..."
	@-$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)

all: VERSION := $(shell . $(BUILD_SUPPORT); getVersion)
all:
	@-echo "Building binaries for all platforms"
	@-echo $(VERSION)
	@. $(MAKE_ALL) ; BIN_PATH=$(BIN_PATH) BINARY_NAME=$(BINARY_NAME) makeall $(VERSION)
