#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = eprinttools

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\`  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

OS = $(shell uname)

EXT = 
ifeq ($(OS), Windows)
        EXT = .exe
endif

QUICK =
ifeq ($(quick), true)
	QUICK = quick=true
endif

PROJECT_LIST = ep eputil doi2eprintxml

build: package $(PROJECT_LIST)

package: eprinttools.go harvest/harvest.go eprint3x.go
	go build

ep: bin/ep$(EXT)

eputil: bin/eputil$(EXT)

doi2eprintxml: bin/doi2eprintxml$(EXT) 

bin/ep$(EXT): eprinttools.go harvest/harvest.go cmd/ep/ep.go
	go build -o bin/ep$(EXT) cmd/ep/ep.go

bin/eputil$(EXT): eprinttools.go harvest/harvest.go eprint3x.go cmd/eputil/eputil.go
	go build -o bin/eputil$(EXT) cmd/eputil/eputil.go

bin/doi2eprintxml$(EXT): eprinttools.go crossref.go datacite.go cmd/doi2eprintxml/doi2eprintxml.go 
	go build -o bin/doi2eprintxml$(EXT) cmd/doi2eprintxml/doi2eprintxml.go

install: 
	env GOBIN=$(GOPATH)/bin go install cmd/ep/ep.go
	env GOBIN=$(GOPATH)/bin go install cmd/eputil/eputil.go
	env GOBIN=$(GOPATH)/bin go install cmd/doi2eprintxml/doi2eprintxml.go


website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css docs/index.md docs/ep.md docs/eputil.md
	./mk-website.bash

test:
	go test
	cd harvest && go test
	cd py && $(MAKE) test $(QUICK)

clean:
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi
	cd py && $(MAKE) clean

dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/ep cmd/ep/ep.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/eputil cmd/eputil/eputil.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md docs/* scripts/* etc/* bin/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/ep.exe cmd/ep/ep.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/eputil.exe cmd/eputil/eputil.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/doi2eprintxml.exe cmd/doi2eprintxml/doi2eprintxml.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md docs/* scripts/* etc/* bin/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/ep cmd/ep/ep.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/eputil cmd/eputil/eputil.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/doi2eprintxml cmd/doi2eprintxml/doi2eprintxml.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md docs/* scripts/* etc/* bin/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/ep cmd/ep/ep.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/eputil cmd/eputil/eputil.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/doi2eprintxml cmd/doi2eprintxml/doi2eprintxml.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md docs/* scripts/* etc/* bin/*
	rm -fR dist/bin
  
distribute_docs:
	mkdir -p dist/etc
	mkdir -p dist/scripts
	mkdir -p dist/docs
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -vR docs/* dist/docs/
	cp -vR scripts/* dist/scripts/
	cp -vR etc/*-example dist/etc/
	./package-versions.bash > dist/package-versions.txt
	#cd py && $(MAKE) release

release: distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish:
	./mk-website.bash
	./publish.bash

