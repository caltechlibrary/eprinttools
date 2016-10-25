#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT_NAME = epgo

build: 
	go build
	go build -o bin/epgo cmds/epgo/epgo.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css
	./mk-website.bash

lint:
	goimports -w epgo.go
	goimports -w epgo_test.go
	goimports -w cmds/epgo/epgo.go
	gofmt -w epgo.go
	gofmt -w epgo_test.go
	gofmt -w cmds/epgo/epgo.go
	golint epgo.go
	golint epgo_test.go
	golint cmds/epgo/epgo.go


test:
	go test

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROJECT_NAME)-binary-release.zip ]; then /bin/rm $(PROJECT_NAME)-binary-release.zip; fi
	if [ -f index.html ]; then /bin/rm *.html; fi

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

