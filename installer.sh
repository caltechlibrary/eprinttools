#!/bin/sh

#
# Set the package name and version to install
#
PACKAGE="eprinttools"
VERSION="v1.3.0"
RELEASE="https://github.com/caltechlibrary/$PACKAGE/releases/tag/$VERSION"

#
# Get the name of this script.
#
INSTALLER=$(basename $0)

#
# Figure out what the zip file is named
#
OS_NAME="$(uname -o)"
MACHINE="$(uname -m)"
case "$OS_NAME" in
   Darwin)
   OS_NAME="macos"
   ;;
esac
ZIPFILE="$PACKAGE-$VERSION-$OS_NAME-$MACHINE.zip"

#
# Check to see if this zip file has been downloaded.
#
DOWNLOAD_URL="https://github.com/caltechlibrary/$PACKAGE/releases/download/$VERSION/$ZIPFILE"
if [ ! -f "$HOME/Downloads/$ZIPFILE" ]; then
	if curl --version > /dev/null 2>&1; then
		curl -L -o "$HOME/Downloads/$ZIPFILE" "$DOWNLOAD_URL"
	else
		cat <<EOT

  To install $PACKAGE you need to download 

    $ZIPFILE

  from 

    $RELEASE

  You can do that with your web browser. After
  that you should be able to re-run $INSTALLER

EOT
		exit 1
	fi
fi

START="$(pwd)"
mkdir -p $HOME/.$PACKAGE/installer
cd $HOME/.$PACKAGE/installer
unzip $HOME/Downloads/$ZIPFILE bin/*

#
# Copy the application into place
#
mkdir -p $HOME/bin
EXPLAIN_OS_POLICY="no"
for APP in $(find bin -type f); do
	FNAME=$(basename $APP)
	mv $APP $HOME/bin/$FNAME
	if ! $HOME/bin/$APP >/dev/null 2>%1; then
		EXPLAIN_OS_POLICY="yes"
	fi
done

#
# Make sure $HOME/bin is in the path
#
DIR_IN_PATH='no'
for P in $PATH; do
  [ "$p" = "$HOME/bin" ] && DIR_IN_PATH='yes'
done
if [ "$DIR_IN_PATH" = "no" ]; then
	echo 'export PATH="$HOME/bin:$PATH"' >>$HOME/.bashrc
	echo 'export PATH="$HOME/bin:$PATH"' >>$HOME/.zshrc
fi
rm -fR $HOME/.$PACKAGE/installer
cd $START

if [ "$EXPLAIN_OS_POLICY" = "no" ]; then
	cat <<EOT

  You need to take additional steps to complete installation.

  Your operating systems' security policied need to "allow"
  running programs from $PACKAGE.

  Example: on macOS you can type open the programs in finder.

      open $HOME/bin

  Find the program(s) and right click and select open
  to enable them to run.

EOT
fi