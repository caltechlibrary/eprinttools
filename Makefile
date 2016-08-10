#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT_NAME = epgo

build: website
	go build
	go build -o bin/epgo cmds/epgo/epgo.go

website: page.tmpl README.md nav.md INSTALL.md LICENSE css/site.css
	./mk-website.bash

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
	./mk-release.bash
	./publish.bash

