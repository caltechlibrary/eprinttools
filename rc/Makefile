#
# Simple Makefile
#
PROJECT = rc

VERSION = $(shell grep 'Version = ' $(PROJECT).go | cut -d \` -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PROG_LIST = apiexplorer

build: $(PROG_LIST)

test:
	go test

lint:
	golint rc.go
	golint cmd/apiexplorer/apiexplorer.go

install:
	env GOBIN=$(GOPATH)/bin go install cmd/apiexplorer/apiexplorer.go


status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

apiexplorer: rc.go cmd/apiexplorer/apiexplorer.go
	env go build -o bin/apiexplorer cmd/apiexplorer/apiexplorer.go


clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -d man ]; then /bin/rm -fR man; fi

man: build
	mkdir -p man/man1
	bin/apiexplorer -generate-manpage | nroff -Tutf8 -man > man/man1/apiexplorer.1

website:
	./mk-website.bash

publish: website
	./publish.bash


dist/linux-amd64:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=amd64 go build -o dist/bin/apiexplorer cmd/apiexplorer/apiexplorer.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md package-versions.txt docs/* etc/* bin/*
	rm -fR dist/bin

dist/windows-amd64:
	mkdir -p dist/bin
	env GOOS=windows GOARCH=amd64 go build -o dist/bin/apiexplorer.exe cmd/apiexplorer/apiexplorer.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md package-versions.txt docs/* etc/* bin/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env GOOS=darwin GOARCH=amd64 go build -o dist/bin/apiexplorer cmd/apiexplorer/apiexplorer.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md package-versions.txt docs/* etc/* bin/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/apiexplorer cmd/apiexplorer/apiexplorer.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md package-versions.txt docs/* etc/* bin/*
	rm -fR dist/bin

distribute_docs:
	mkdir -p dist
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	./package-versions.bash > dist/package-versions.txt
	if [ -d docs ]; then mkdir -p dist/docs; cp -v docs/*.md dist/docs/; fi
	if [ -d etc ]; then mkdir -p dist/etc; cp -v etc/*-example dist/etc/; fi

release: distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7

