#
# Simple Makefile for conviently testing, building and deploying experiment.
#
build: *.go
	go build
	go build -o bin/epgo cmds/epgo/epgo.go

test:
	go test

clean:
	if [ -d bin ]; then rm bin/*; fi
