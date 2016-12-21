#!/bin/bash
#
# Make releases for Linux/amd64, Linux/ARM6 and Linux/ARM7 (Raspberry Pi), Windows, and Mac OX X (darwin)
#
PROJECT=epgo
VERSION=$(grep -m 1 'Version =' $PROJECT.go | cut -d\" -f 2)
RELEASE_NAME=$PROJECT-$VERSION
echo "Preparing release $RELEASE_NAME"
for PROGNAME in epgo genpages indexpages sitemapper servepages; do
  echo "Cross compiling $PROGNAME"
  env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
  env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
  env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspberrypi-arm6/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
  env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
  env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/$PROGNAME.exe cmds/$PROGNAME/$PROGNAME.go
done

mkdir -p dist/etc/systemd/system
mkdir -p dist/htdocs/css
mkdir -p dist/htdocs/js
mkdir -p dist/htdocs/assets
for FNAME in README.md LICENSE INSTALL.md NOTES.md templates scripts; do
  cp -vR $FNAME dist/
done
cp -vR etc/*-example dist/etc/
cp -vR etc/systemd/system/*-example dist/etc/systemd/system/
cp -vR htdocs/index.* dist/htdocs/
cp -vR htdocs/css dist/htdocs/
cp -vR htdocs/js dist/htdocs/
cp -vR htdocs/assets dist/htdocs/
echo "Zipping $RELEASE_NAME-release.zip"
zip -r "$RELEASE_NAME-release.zip" dist/*
