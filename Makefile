#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = eprinttools

GIT_GROUP = caltechlibrary

RELEASE_DATE=$(shell date +'%Y-%m-%d')

RELEASE_HASH=$(shell git log --pretty=format:'%h' -n 1)

VERSION = $(shell grep '"version":' codemeta.json | cut -d\"  -f 4)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)


MAN_PAGES = doi2eprintxml.1 eputil.1 ep3apid.1 epfmt.1 ep3harvester.1 ep3genfeeds.1 ep3datasets.1

PROGRAMS = $(shell ls -1 cmd)

PACKAGE = $(shell ls -1 *.go cleaner/*.go clsrules/*.go)

PANDOC = $(shell which pandoc)

OS = $(shell uname)

#PREFIX = /usr/local/bin
PREFIX = $(HOME)

ifneq ($(prefix),)
	PREFIX = $(prefix)
endif

EXT = 
ifeq ($(OS), Windows)
        EXT = .exe
endif

QUICK =
ifeq ($(quick), true)
	QUICK = quick=true
endif


build: version.go $(PROGRAMS) about.md installer.sh

version.go: .FORCE
	echo '' | pandoc --from t2t --to plain \
                --metadata-file codemeta.json \
                --metadata package=$(PROJECT) \
                --metadata version=$(VERSION) \
                --metadata release_date=$(RELEASE_DATE) \
                --metadata release_hash=$(RELEASE_HASH) \
                --template codemeta-version-go.tmpl \
                LICENSE >version.go

$(PROGRAMS): $(PACKAGE)
	@mkdir -p bin
	go build -o bin/$@$(EXT) cmd/$@/$@.go

man: $(MAN_PAGES)

$(MAN_PAGES): .FORCE
	mkdir -p man/man1
	$(PANDOC) $@.md --from markdown --to man -s >man/man1/$@

CITATION.cff: .FORCE
	cat codemeta.json | sed -E   's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	if [ -f $(PANDOC) ]; then echo "" | $(PANDOC) --metadata title="Cite $(PROJECT)" --metadata-file=_codemeta.json --template=codemeta-cff.tmpl >CITATION.cff; fi
	if [ -f _codemeta.json ]; then rm _codemeta.json; fi

# NOTE: macOS requires a "mv" command for placing binaries instead of "cp" due to signing process of compile
install: build man
	@echo "Installing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f ./bin/$$FNAME ]; then mv ./bin/$$FNAME $(PREFIX)/bin/$$FNAME; fi; done
	@mkdir -p $(PREFIX)/man/man1
	@for FNAME in $(MAN_PAGES); do if [ -f ./man/man1/$$FNAME ]; then cp -v ./man/man1/$$FNAME $(PREFIX)/man/man1/$$FNAME; fi; done
	@echo ""
	@echo "Make sure $(PREFIX)/bin is in your PATH"


uninstall: .FORCE
	@echo "Removing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f $(PREFIX)/bin/$$FNAME ]; then rm -v $(PREFIX)/bin/$$FNAME; fi; done
	@for FNAME in $(MAN_PAGES); do if [ -f $(PREFIX)/man/man1/$$FNAME ]; then rm -v $(PREFIX)/man/man1/$$FNAME; fi; done

index.md: .FORCE
	cp README.md index.md
	git add index.md

about.md: .FORCE
	cat codemeta.json | sed -E 's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	if [ -f $(PANDOC) ]; then echo "" | pandoc --metadata title="About $(PROJECT)" --metadata-file=_codemeta.json --template codemeta-about.tmpl >about.md; fi
	if [ -f _codemeta.json ]; then rm _codemeta.json; fi

installer.sh: .FORCE
	@echo '' | pandoc --metadata title="Installer" --metadata git_org_or_person="$(GIT_GROUP)" --metadata-file codemeta.json --template codemeta-installer.tmpl >installer.sh
	chmod 775 installer.sh
	git add -f installer.sh

website: index.md about.md page.tmpl *.md LICENSE css/site.css
	make -f website.mak


test: version.go eputil epfmt doi2eprintxml ep3apid
	- cd cleaner && go test -test.v
	- cd clsrules && go test -test.v
	- go test -timeout 1h -test.v
	./test_cmds.bash


clean:
	@if [ -f version.go ]; then rm version.go; fi
	@if [ -d bin ]; then rm -fR bin; fi
	@if [ -d dist ]; then rm -fR dist; fi

dist/Linux-x86_64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env  GOOS=linux GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/$$FNAME.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-x86_64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/* 
	@rm -fR dist/bin

dist/macOS-x86_64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=amd64 go build -o dist/bin/$$FNAME cmd/$$FNAME/$$FNAME.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macOS-x86_64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin

dist/macOS-arm64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=arm64 go build -o dist/bin/$$FNAME cmd/$$FNAME/$$FNAME.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-macOS-arm64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin

dist/Windows-x86_64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=amd64 go build -o dist/bin/$$FNAME.exe cmd/$$FNAME/$$FNAME.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Windows-x86_64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin

dist/Windows-arm64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=arm64 go build -o dist/bin/$$FNAME.exe cmd/$$FNAME/$$FNAME.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Windows-arm64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin

# Raspberry Pi running Ubuntu 64bit or Apple M1 running a VM
dist/Linux-aarch64:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=arm64 go build -o dist/bin/$$FNAME cmd/$$FNAME/$$FNAME.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-aarch64.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin
  

# Raspberry Pi OS (32bit) as reported on Raspberry Pi Model 3B+
dist/Linux-armv7l:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/$$FNAME cmd/$$FNAME/$$FNAME.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-Linux-armv7l.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin
  
dist/RaspberryPiOS-arm7:
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/$$FNAME cmd/$$FNAME/$$FNAME.go; done
	@cd dist && zip -r $(PROJECT)-v$(VERSION)-RaspberryPiOS-arm7.zip LICENSE codemeta.json CITATION.cff *.md bin/* man/*
	@rm -fR dist/bin
  
distribute_docs:
	mkdir -p dist/man/man1
	cp -v codemeta.json dist/
	cp -v CITATION.cff dist/
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v installer.sh dist/
	cp -vR man dist/

dist_clean: .FORCE
	if [ -d dist ]; then rm -fR dist/*; fi

release: build man CITATION.cff dist_cleanup distribute_docs dist/Linux-x86_64 dist/Windows-x86_64 dist/Windows-arm64 dist/macOS-x86_64 dist/macOS-arm64 dist/RaspberryPiOS-arm7 dist/Linux-aarch64 dist/Linux-armv7l

status:
	git status

save:
	@if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

publish: website
	./publish.bash

.FORCE:
