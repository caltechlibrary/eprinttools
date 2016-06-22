#
# Simple Makefile for conviently testing, building and deploying experiment.
#
build:
	go build
	go build -o bin/epgo cmds/epgo/epgo.go

test:
	go test

clean:
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi

release:
	./mk-release.sh

