#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = epgo

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\"  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PROJECT_LIST = epgo epgo-genpages

build: package $(PROJECT_LIST)

package: epgo.go
	go build

epgo: bin/epgo

bin/epgo: epgo.go  harvest.go grantNumbers.go funders.go cmds/epgo/epgo.go
	go build -o bin/epgo cmds/epgo/epgo.go

epgo-genpages: bin/epgo-genpages

bin/epgo-genpages: epgo.go  harvest.go grantNumbers.go funders.go cmds/epgo-genpages/epgo-genpages.go
	go build -o bin/epgo-genpages cmds/epgo-genpages/epgo-genpages.go

install: 
	env GOBIN=$(HOME)/bin go install cmds/epgo/epgo.go
	env GOBIN=$(HOME)/bin go install cmds/epgo-genpages/epgo-genpages.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css htdocs/index.md
	mkpage "content=htdocs/index.md" templates/default/index.html > htdocs/index.html
	./mk-website.bash

format:
	gofmt -w epgo.go
	gofmt -w epgo_test.go
	gofmt -w harvest.go
	gofmt -w funders.go
	gofmt -w grantNumbers.go
	gofmt -w cmds/epgo/epgo.go
	gofmt -w cmds/epgo-genpages/epgo-genpages.go

lint:
	golint epgo.go
	golint epgo_test.go
	golint harvest.go
	golint funders.go
	golint grantNumbers.go
	golint cmds/epgo/epgo.go
	golint cmds/epgo-genpages/epgo-genpages.go

test:
	go test

clean:
	if [ -f index.html ]; then /bin/rm *.html; fi
	if [ -d htdocs/person ]; then /bin/rm -fR htdocs/person; fi
	if [ -d htdocs/affiliation ]; then /bin/rm -fR htdocs/affiliation; fi
	if [ -d htdocs/recent ]; then /bin/rm -fR htdocs/recent; fi
	if [ -d htdocs/repository ]; then /bin/rm -fR htdocs/repository; fi
	if [ -d htdocs/funder ]; then /bin/rm -fR htdocs/funder; fi
	if [ -d htdocs/grantNumber ]; then /bin/rm -fR htdocs/grantNumber; fi
	if [ "$(EPGO_REPOSITORY_PATH)" != "" ] && [ -d htdocs/$(EPGO_REPOSITORY_PATH) ]; then /bin/rm -fR htdocs/$(EPGO_REPOSITORY_PATH); fi
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROJECT)-$(VERSION)-release.zip ]; then /bin/rm $(PROJECT)-$(VERSION)-release.zip; fi

dist/linux-amd64:
	env  GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/epgo cmds/epgo/epgo.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/epgo-genpages cmds/epgo-genpages/epgo-genpages.go

dist/windows-amd64:
	env  GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/epgo.exe cmds/epgo/epgo.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/epgo-genpages.exe cmds/epgo-genpages/epgo-genpages.go

dist/macosx-amd64:
	env  GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/epgo cmds/epgo/epgo.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/epgo-genpages cmds/epgo-genpages/epgo-genpages.go

dist/raspbian-arm7:
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/epgo cmds/epgo/epgo.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/epgo-genpages cmds/epgo-genpages/epgo-genpages.go
  
release: dist/linux-amd64 dist/windows-amd64 macosx-amd64 raspbian-arm7
	mkdir -p dist/etc/systemd/system
	mkdir -p dist/htdocs/css
	mkdir -p dist/htdocs/js
	mkdir -p dist/htdocs/assets
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -vR templates dist/
	cp -vR scripts dist/
	cp -vR etc/*-example dist/etc/
	cp -vR etc/systemd/system/*-example dist/etc/systemd/system/
	cp -vR htdocs/index.* dist/htdocs/
	cp -vR htdocs/css dist/htdocs/
	cp -vR htdocs/js dist/htdocs/
	cp -vR htdocs/assets dist/htdocs/
	zip -r $(PROJECT)-$(VERSION)-release.zip dist/*

status:
	git status

save:
	git commit -am "Quick Save"
	git push origin $(BRANCH)

publish:
	./mk-website.bash
	./publish.bash

