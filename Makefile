#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = ep

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\"  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PROJECT_LIST = ep

build: package $(PROJECT_LIST)

package: ep.go
	go build

ep: bin/ep

bin/ep: ep.go  harvest.go grantNumbers.go funders.go cmds/ep/ep.go
	go build -o bin/ep cmds/ep/ep.go

bin/ep-genfeeds: ep.go harvest.go grantNumnberrs.go funders.go cmds/ep-genfeeds/ep-genfeeds.go
	go build -o bin/ep-genfeeds cmds/ep-genfeeds.go

install: 
	env GOBIN=$(HOME)/bin go install cmds/ep/ep.go
	env GOBIN=$(HOME)/bin go install cmds/ep-genfeeds/ep-genfeeds.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css
	./mk-website.bash

format:
	gofmt -w ep.go
	gofmt -w ep_test.go
	gofmt -w harvest.go
	gofmt -w funders.go
	gofmt -w grantNumbers.go
	gofmt -w cmds/ep/ep.go

lint:
	golint ep.go
	golint ep_test.go
	golint harvest.go
	golint funders.go
	golint grantNumbers.go
	golint cmds/ep/ep.go

test:
	go test

clean:
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi

dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/ep cmds/ep/ep.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md docs/* scripts/* etc/* bin/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/ep.exe cmds/ep/ep.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md docs/* scripts/* etc/* bin/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/ep cmds/ep/ep.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md docs/* scripts/* etc/* bin/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/ep cmds/ep/ep.go
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

release: distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish:
	./mk-website.bash
	./publish.bash

