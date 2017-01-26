#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = epgo

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\"  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PROJECT_LIST = epgo epgo-genpages epgo-indexpages epgo-sitemapper epgo-servepages

build: package $(PROJECT_LIST)

package: epgo.go
	go build

epgo: bin/epgo

bin/epgo: epgo.go api.go export.go cmds/epgo/epgo.go
	go build -o bin/epgo cmds/epgo/epgo.go

epgo-genpages: bin/epgo-genpages

bin/epgo-genpages: epgo.go api.go export.go cmds/epgo-genpages/epgo-genpages.go
	go build -o bin/epgo-genpages cmds/epgo-genpages/epgo-genpages.go

epgo-indexpages: bin/epgo-indexpages

bin/epgo-indexpages: epgo.go api.go export.go cmds/epgo-indexpages/epgo-indexpages.go
	go build -o bin/epgo-indexpages cmds/epgo-indexpages/epgo-indexpages.go

epgo-sitemapper: bin/epgo-sitemapper

bin/epgo-sitemapper: epgo.go api.go export.go cmds/epgo-sitemapper/epgo-sitemapper.go
	go build -o bin/epgo-sitemapper cmds/epgo-sitemapper/epgo-sitemapper.go

epgo-servepages: bin/epgo-servepages

bin/epgo-servepages: epgo.go api.go export.go cmds/epgo-servepages/epgo-servepages.go
	go build -o bin/epgo-servepages cmds/epgo-servepages/epgo-servepages.go
	mkpage "content=htdocs/index.md" templates/default/index.html > htdocs/index.html

install: $(PROJECT_LIST)
	env GOBIN=$(HOME)/bin go install cmds/epgo/epgo.go
	env GOBIN=$(HOME)/bin go install cmds/epgo-genpages/epgo-genpages.go
	env GOBIN=$(HOME)/bin go install cmds/epgo-indexpages/epgo-indexpages.go
	env GOBIN=$(HOME)/bin go install cmds/epgo-sitemapper/epgo-sitemapper.go
	env GOBIN=$(HOME)/bin go install cmds/epgo-servepages/epgo-servepages.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css htdocs/index.md
	mkpage "content=htdocs/index.md" templates/default/index.html > htdocs/index.html
	./mk-website.bash

format:
	goimports -w epgo.go
	goimports -w epgo_test.go
	goimports -w cmds/epgo/epgo.go
	goimports -w cmds/epgo-genpages/epgo-genpages.go
	goimports -w cmds/epgo-indexpages/epgo-indexpages.go
	goimports -w cmds/epgo-sitemapper/epgo-sitemapper.go
	goimports -w cmds/epgo-servepages/epgo-servepages.go
	gofmt -w epgo.go
	gofmt -w epgo_test.go
	gofmt -w cmds/epgo/epgo.go
	gofmt -w cmds/epgo-genpages/epgo-genpages.go
	gofmt -w cmds/epgo-indexpages/epgo-indexpages.go
	gofmt -w cmds/epgo-sitemapper/epgo-sitemapper.go
	gofmt -w cmds/epgo-servepages/epgo-servepages.go

lint:
	golint epgo.go
	golint epgo_test.go
	golint cmds/epgo/epgo.go
	golint cmds/epgo-genpages/epgo-genpages.go
	golint cmds/epgo-indexpages/epgo-indexpages.go
	golint cmds/epgo-sitemapper/epgo-sitemapper.go
	golint cmds/epgo-servepages/epgo-servepages.go


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
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/epgo cmds/epgo/epgo.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/epgo-genpages cmds/epgo-genpages/epgo-genpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/epgo-indexpages cmds/epgo-indexpages/epgo-indexpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/epgo-servepages cmds/epgo-servepages/epgo-servepages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/epgo-sitemapper cmds/epgo-sitemapper/epgo-sitemapper.go

dist/windows-amd64:
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/epgo.exe cmds/epgo/epgo.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/epgo-genpages.exe cmds/epgo-genpages/epgo-genpages.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/epgo-indexpages.exe cmds/epgo-indexpages/epgo-indexpages.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/epgo-servepages.exe cmds/epgo-servepages/epgo-servepages.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/epgo-sitemapper.exe cmds/epgo-sitemapper/epgo-sitemapper.go

dist/macosx-amd64:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/epgo cmds/epgo/epgo.go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/epgo-genpages cmds/epgo-genpages/epgo-genpages.go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/epgo-indexpages cmds/epgo-indexpages/epgo-indexpages.go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/epgo-servepages cmds/epgo-servepages/epgo-servepages.go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/epgo-sitemapper cmds/epgo-sitemapper/epgo-sitemapper.go

dist/raspbian-arm7:
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/epgo cmds/epgo/epgo.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/epgo-genpages cmds/epgo-genpages/epgo-genpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/epgo-indexpages cmds/epgo-indexpages/epgo-indexpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/epgo-servepages cmds/epgo-servepages/epgo-servepages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/epgo-sitemapper cmds/epgo-sitemapper/epgo-sitemapper.go
  
dist/raspbian-arm6:
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspberrypi-arm6/epgo cmds/epgo/epgo.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspberrypi-arm6/epgo-genpages cmds/epgo-genpages/epgo-genpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspberrypi-arm6/epgo-indexpages cmds/epgo-indexpages/epgo-indexpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspberrypi-arm6/epgo-servepages cmds/epgo-servepages/epgo-servepages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspberrypi-arm6/epgo-sitemapper cmds/epgo-sitemapper/epgo-sitemapper.go


release: dist/linux-amd64 dist/windows-amd64 macosx-amd64 raspbian-arm7 raspbian-arm6
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
	git commit -am "Quick save"
	git push origin $(BRANCH)

publish:
	./mk-website.bash
	./publish.bash

