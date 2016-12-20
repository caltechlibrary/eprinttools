#!/bin/bash
#
# Make releases for Linux/amd64, Linux/ARM6 and Linux/ARM7 (Raspberry Pi), Windows, and Mac OX X (darwin)
#
PROJECT=epgo
VERSION=$(grep -m 1 'Version =' epgo.go | cut -d\" -f 2)
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
echo "Zipping $RELEASE_NAME-release.zip"
zip -r "$RELEASE_NAME-release.zip" README.md INSTALL.md LICENSE scripts/* etc/*example etc/systemd/system/*example templates/* htdocs/index.* htdocs/css/* htdocs/js/* htdocs/assets/* dist/*
