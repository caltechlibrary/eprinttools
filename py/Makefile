#
# Simple Makefile for building C-Shared library and Python3 module.
#
PROJECT = eprinttools

LIB_NAME = libeprint3

VERSION = $(shell grep -m 1 'Version =' ../$(PROJECT).go | cut -d\`  -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

ARCH = x86_64

OS = $(shell uname)

EXT = .so
ifeq ($(OS), Windows)
	EXT = .dll
	ARCH = x86_64
	os = windows
endif
ifeq ($(OS), Darwin)
	EXT = .dylib
	ARCH = $(shell arch)
	OS = macosx
endif
ifeq ($(OS), Linux)
	EXT = .so
	ARCH = $(shell arch)
	OS = linux
endif
ifeq ($(ARCH), i386)
	ARCH = amd64
endif
ifeq ($(ARCH), x86_64)
	ARCH = amd64
endif

build: $(PROJECT)/$(LIB_NAME)$(EXT)
	python3 setup.py build

$(PROJECT)/$(LIB_NAME)$(EXT): $(LIB_NAME).go
	go build -buildmode=c-shared -o "$(PROJECT)/$(LIB_NAME)$(EXT)" "$(LIB_NAME).go"

clean:
	if [ -f "MANIFEST" ]; then rm MANIFEST; fi
	if [ -f "$(PROJECT)/$(LIB_NAME)$(EXT)" ]; then rm "$(PROJECT)/$(LIB_NAME)$(EXT)"; fi
	if [ -f "$(PROJECT)/$(LIB_NAME).h" ]; then rm "$(PROJECT)/$(LIB_NAME).h"; fi
	if [ -d "dist" ]; then rm -fR dist; fi
	if [ -d "build" ]; then rm -fR build; fi
	if [ -d "__pycache__" ]; then rm -fR "__pycache__"; fi	
	if [ -d "$(PROJECT)/__pycache__" ]; then rm -fR "$(PROJECT)/__pycache__"; fi	

test: clean $(PROJECT)/$(LIB_NAME)$(EXT)
	python3 $(PROJECT)_test.py
	python3 $(PROJECT)_dataset_test.py

release: clean $(PROJECT)/$(LIB_NAME)$(EXT)
	go build -buildmode=c-shared -o "$(PROJECT)/$(LIB_NAME)$(EXT)" "$(LIB_NAME).go"
	python3 setup.py sdist
	mkdir -p ../dist
	cp -v dist/$(PROJECT)-$(VERSION).tar.gz "../dist/py3-$(PROJECT)-$(VERSION)-$(OS)-$(ARCH).tar.gz"

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

