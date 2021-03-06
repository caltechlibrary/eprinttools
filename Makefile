#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = eprinttools

VERSION = $(shell grep -m 1 'Version =' version.go | cut -d\`  -f 2)

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

PROJECT_LIST = eputil epfmt doi2eprintxml eprintxml2json

build: package $(PROJECT_LIST)

package: eprinttools.go eprint3x.go
	go build


eputil: bin/eputil$(EXT)

epfmt: bin/epfmt$(EXT)

doi2eprintxml: bin/doi2eprintxml$(EXT) 

eprintxml2json: bin/eprintxml2json$(EXT)

bin/eputil$(EXT): eprinttools.go eprint3x.go cmd/eputil/eputil.go
	go build -o bin/eputil$(EXT) cmd/eputil/eputil.go

bin/epfmt$(EXT): eprinttools.go eprint3x.go cmd/epfmt/epfmt.go
	go build -o bin/epfmt$(EXT) cmd/epfmt/epfmt.go

bin/doi2eprintxml$(EXT): eprinttools.go crossref.go datacite.go clsrules/clsrules.go cmd/doi2eprintxml/doi2eprintxml.go 
	go build -o bin/doi2eprintxml$(EXT) cmd/doi2eprintxml/doi2eprintxml.go

bin/eprintxml2json$(EXT): eprinttools.go eprint3x.go cmd/eprintxml2json/eprintxml2json.go 
	go build -o bin/eprintxml2json$(EXT) cmd/eprintxml2json/eprintxml2json.go

install: 
	env GOBIN=$(GOPATH)/bin go install cmd/eputil/eputil.go
	env GOBIN=$(GOPATH)/bin go install cmd/epfmt/epfmt.go
	env GOBIN=$(GOPATH)/bin go install cmd/doi2eprintxml/doi2eprintxml.go
	env GOBIN=$(GOPATH)/bin go install cmd/eprintxml2json/eprintxml2json.go


website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css docs/index.md docs/eputil.md
	./mk-website.bash



test: eputil epfmt doi2eprintxml eprintxml2json
	go test -timeout 45m
	./test_cmds.bash

clean:
	if [ -d htdocs ]; then rm -fR htdocs; fi
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi
	if [ -d man ]; then rm -fR man; fi

man: build
	mkdir -p man/man1
	bin/ep -generate-manpage | nroff -Tutf8 -man > man/man1/ep.1
	bin/eputil -generate-manpage | nroff -Tutf8 -man > man/man1/eputil.1
	bin/epfmt -generate-manpage | nroff -Tutf8 -man > man/man1/epfmt.1
	bin/doi2eprintxml -generate-manpage | nroff -Tutf8 -man > man/man1/doi2eprintxml.1
	bin/eprintxml2json -generate-manpage | nroff -Tutf8 -man > man/man1/eprintxml2json.1

dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/eputil cmd/eputil/eputil.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/epfmt cmd/epfmt/epfmt.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/doi2eprintxml cmd/doi2eprintxml/doi2eprintxml.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/eprintxml2json cmd/eprintxml2json/eprintxml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md docs/* eprinttools/* bin/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/eputil.exe cmd/eputil/eputil.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/epfmt.exe cmd/epfmt/epfmt.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/doi2eprintxml.exe cmd/doi2eprintxml/doi2eprintxml.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/eprintxml2json.exe cmd/eprintxml2json/eprintxml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md docs/* eprinttools/* bin/*
	rm -fR dist/bin

dist/macos-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/eputil cmd/eputil/eputil.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/epfmt cmd/epfmt/epfmt.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/doi2eprintxml cmd/doi2eprintxml/doi2eprintxml.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/eprintxml2json cmd/eprintxml2json/eprintxml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macos-amd64.zip README.md LICENSE INSTALL.md docs/* eprinttools/* bin/*
	rm -fR dist/bin

dist/macos-arm64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/eputil cmd/eputil/eputil.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/epfmt cmd/epfmt/epfmt.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/doi2eprintxml cmd/doi2eprintxml/doi2eprintxml.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/eprintxml2json cmd/eprintxml2json/eprintxml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macos-arm64.zip README.md LICENSE INSTALL.md docs/* eprinttools/* bin/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/eputil cmd/eputil/eputil.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/epfmt cmd/epfmt/epfmt.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/doi2eprintxml cmd/doi2eprintxml/doi2eprintxml.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/eprintxml2json cmd/eprintxml2json/eprintxml2json.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md docs/* eprinttools/* bin/*
	rm -fR dist/bin
  
distribute_python:
	mkdir -p dist/eprinttools/eprints3x
	mkdir -p dist/eprinttools/eprintviews
	cp -v eprinttools/eprints3x/*.py dist/eprinttools/eprints3x/
	cp -vR eprinttools/eprintviews/*.py dist/eprinttools/eprintviews/
	cp -vR static dist/
	cp -vR templates dist/
	cp config.json-example dist/
	cp setup.py dist/
	cp harvester_full.py dist/
	cp harvester_recent.py dist/
	cp genviews.py dist/
	cp indexer.py dist/
	cp mk_website.py dist/
	cp publisher.py dist/
	cp invalidate_cloudfront.py dist/

distribute_docs:
	mkdir -p dist/docs
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -vR docs/* dist/docs/

release: distribute_docs distribute_python dist/linux-amd64 dist/windows-amd64 dist/macos-amd64 dist/macos-arm64 dist/raspbian-arm7

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish: website
	./publish.bash

