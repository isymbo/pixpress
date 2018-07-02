LDFLAGS += -X "github.com/isymbo/pixpress/setting.BuildTime=$(shell date '+%Y-%m-%d %H:%M:%S %Z')"
LDFLAGS += -X "github.com/isymbo/pixpress/setting.BuildGitHash=$(shell git rev-parse HEAD)"

DATA_FILES := $(shell find config | sed 's/ /\\ /g')
LESS_FILES := $(wildcard public/less/pixpress.less public/less/_*.less)
# GENERATED  := pkg/bindata/bindata.go public/css/pixpress.css

OS := $(shell uname)

TAGS = ""
BUILD_FLAGS = "-v"

RELEASE_ROOT = "release"
RELEASE_PIXPRESS = "release/pixpress"
NOW = $(shell date '+%Y%m%d%H%M%S_%Z')
GOVET = go tool vet -composites=false -methods=false -structtags=false
GOPATH = $(shell echo $${PWD%%src*})

.PHONY: build pack release bindata clean

.IGNORE: public/css/pixpress.css

all: build

check: test

dist: release

web: build
	./pixpress web

govet:
	$(GOVET) pixpress.go
	$(GOVET) app cmd setting util

build: $(GENERATED)
	go install $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -tags '$(TAGS)'
	cp '$(GOPATH)/bin/pixpress' .

build-dev: $(GENERATED) govet
	go install $(BUILD_FLAGS) -tags '$(TAGS)'
	cp '$(GOPATH)/bin/pixpress' .

build-dev-race: $(GENERATED) govet
	go install $(BUILD_FLAGS) -race -tags '$(TAGS)'
	cp '$(GOPATH)/bin/pixpress' .

pack:
	rm -rf $(RELEASE_PIXPRESS)
	mkdir -p $(RELEASE_PIXPRESS)
	cp -r pixpress LICENSE README.md README_ZH.md templates public scripts $(RELEASE_PIXPRESS)
	rm -rf $(RELEASE_PIXPRESS)/public/config.codekit $(RELEASE_PIXPRESS)/public/less
	cd $(RELEASE_ROOT) && zip -r pixpress.$(NOW).zip "pixpress"

release: build pack

bindata: pkg/bindata/bindata.go

pkg/bindata/bindata.go: $(DATA_FILES)
	go-bindata -o=$@ -ignore="\\.DS_Store|README.md|TRANSLATORS|auth.d" -pkg=bindata conf/...

less: public/css/pixpress.css

public/css/pixpress.css: $(LESS_FILES)
	lessc $< $@

clean:
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -print0 | xargs -0 rm

test:
	go test -cover -race ./...

fixme:
	grep -rnw "FIXME" cmd app setting util

todo:
	grep -rnw "TODO" cmd app setting util

cloc:
	gocloc .

# Legacy code should be remove by the time of release
legacy:
	grep -rnw "LEGACY" cmd app setting util
