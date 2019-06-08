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

PROJECT_LIST = ep eputil epfmt doi2eprintxml eprintxml2json

build: package $(PROJECT_LIST)

package: eprinttools.go harvest/harvest.go eprint3x.go
	go build

ep: bin/ep$(EXT)

eputil: bin/eputil$(EXT)

epfmt: bin/epfmt$(EXT)

doi2eprintxml: bin/doi2eprintxml$(EXT) 

eprintxml2json: bin/eprintxml2json$(EXT)

bin/ep$(EXT): eprinttools.go harvest/harvest.go cmd/ep/ep.go
	go build -o bin/ep$(EXT) cmd/ep/ep.go

bin/eputil$(EXT): eprinttools.go harvest/harvest.go eprint3x.go cmd/eputil/eputil.go
	go build -o bin/eputil$(EXT) cmd/eputil/eputil.go

bin/epfmt$(EXT): eprinttools.go harvest/harvest.go eprint3x.go cmd/epfmt/epfmt.go
	go build -o bin/epfmt$(EXT) cmd/epfmt/epfmt.go

bin/doi2eprintxml$(EXT): eprinttools.go crossref.go datacite.go clsrules/clsrules.go cmd/doi2eprintxml/doi2eprintxml.go 
	go build -o bin/doi2eprintxml$(EXT) cmd/doi2eprintxml/doi2eprintxml.go

bin/eprintxml2json$(EXT): eprinttools.go eprint3x.go cmd/eprintxml2json/eprintxml2json.go 
	go build -o bin/eprintxml2json$(EXT) cmd/eprintxml2json/eprintxml2json.go

install: 
	env GOBIN=$(GOPATH)/bin go install cmd/ep/ep.go
	env GOBIN=$(GOPATH)/bin go install cmd/eputil/eputil.go
	env GOBIN=$(GOPATH)/bin go install cmd/epfmt/epfmt.go
	env GOBIN=$(GOPATH)/bin go install cmd/doi2eprintxml/doi2eprintxml.go
	env GOBIN=$(GOPATH)/bin go install cmd/eprintxml2json/eprintxml2json.go


website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css docs/index.md docs/ep.md docs/eputil.md
	./mk-website.bash

test:
	go test -timeout 45m
	cd harvest && go test
	./test_cmds.bash

clean:
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi
	if [ -d man ]; then rm -fR man; fi
	#cd py && $(MAKE) clean

man: build
	mkdir -p man/man1
	bin/ep -generate-manpage | nroff -Tutf8 -man > man/man1/ep.1
	bin/eputil -generate-manpage | nroff -Tutf8 -man > man/man1/eputil.1
	bin/epfmt -generate-manpage | nroff -Tutf8 -man > man/man1/epfmt.1
	bin/doi2eprintxml -generate-manpage | nroff -Tutf8 -man > man/man1/doi2eprintxml.1
	bin/eprintxml2json -generate-manpage | nroff -Tutf8 -man > man/man1/eprintxml2json.1

dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/ep cmd/ep/ep.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/eputil cmd/eputil/eputil.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/epfmt cmd/epfmt/epfmt.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/doi2eprintxml cmd/doi2eprintxml/doi2eprintxml.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/eprintxml2json cmd/eprintxml2json/eprintxml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md docs/* scripts/* etc/* bin/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/ep.exe cmd/ep/ep.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/eputil.exe cmd/eputil/eputil.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/epfmt.exe cmd/epfmt/epfmt.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/doi2eprintxml.exe cmd/doi2eprintxml/doi2eprintxml.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/eprintxml2json.exe cmd/eprintxml2json/eprintxml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md docs/* scripts/* etc/* bin/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/ep cmd/ep/ep.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/eputil cmd/eputil/eputil.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/epfmt cmd/epfmt/epfmt.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/doi2eprintxml cmd/doi2eprintxml/doi2eprintxml.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/eprintxml2json cmd/eprintxml2json/eprintxml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md docs/* scripts/* etc/* bin/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/ep cmd/ep/ep.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/eputil cmd/eputil/eputil.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/epfmt cmd/epfmt/epfmt.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/doi2eprintxml cmd/doi2eprintxml/doi2eprintxml.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/eprintxml2json cmd/eprintxml2json/eprintxml2json.go
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

