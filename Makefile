#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT_NAME = epgo

PROJECT_LIST = epgo genpages indexpages sitemapper servepages

build: package $(PROJECT_LIST)

package: epgo.go
	go build

epgo: epgo.go cmds/epgo/epgo.go
	go build -o bin/epgo cmds/epgo/epgo.go

genpages: epgo.go cmds/genpages/genpages.go
	go build -o bin/genpages cmds/genpages/genpages.go

indexpages: epgo.go cmds/indexpages/indexpages.go
	go build -o bin/indexpages cmds/indexpages/indexpages.go

sitemapper: epgo.go cmds/sitemapper/sitemapper.go
	go build -o bin/sitemapper cmds/sitemapper/sitemapper.go

servepages: epgo.go cmds/servepages/servepages.go
	go build -o bin/servepages cmds/servepages/servepages.go
	mkpage "content=htdocs/index.md" page.tmpl > htdocs/index.html

install: $(PROJECT_LIST)
	go install cmds/epgo/epgo.go
	go install cmds/genpages/genpages.go
	go install cmds/indexpages/indexpages.go
	go install cmds/sitemapper/sitemapper.go
	go install cmds/servepages/servepages.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css htdocs/index.md
	mkpage "content=htdocs/index.md" templates/default/index.html > htdocs/index.html
	./mk-website.bash

format:
	goimports -w epgo.go
	goimports -w epgo_test.go
	goimports -w cmds/epgo/epgo.go
	goimports -w cmds/genpages/genpages.go
	goimports -w cmds/indexpages/indexpages.go
	goimports -w cmds/sitemapper/sitemapper.go
	goimports -w cmds/servepages/servepages.go
	gofmt -w epgo.go
	gofmt -w epgo_test.go
	gofmt -w cmds/epgo/epgo.go
	gofmt -w cmds/genpages/genpages.go
	gofmt -w cmds/indexpages/indexpages.go
	gofmt -w cmds/sitemapper/sitemapper.go
	gofmt -w cmds/servepages/servepages.go

lint:
	golint epgo.go
	golint epgo_test.go
	golint cmds/epgo/epgo.go
	golint cmds/genpages/genpages.go
	golint cmds/indexpages/indexpages.go
	golint cmds/sitemapper/sitemapper.go
	golint cmds/servepages/servepages.go


test:
	go test

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROJECT_NAME)-release.zip ]; then /bin/rm $(PROJECT_NAME)-release.zip; fi
	if [ -f index.html ]; then /bin/rm *.html; fi
	if [ -d htdocs/person ]; then /bin/rm -fR htdocs/person; fi
	if [ -d htdocs/affiliation ]; then /bin/rm -fR htdocs/affiliation; fi
	if [ -d htdocs/recent ]; then /bin/rm -fR htdocs/recent; fi
	if [ -d htdocs/repository ]; then /bin/rm -fR htdocs/repository; fi
	if [ -d htdocs/$(EPGO_REPOSITORY_PATH) ]; then /bin/rm -fR htdocs/$(EPGO_REPOSITORY_PATH); fi

release:
	./mk-release.bash

status:
	git status

save:
	git commit -am "Quick save"
	git push origin master

publish:
	./mk-website.bash
	./publish.bash

