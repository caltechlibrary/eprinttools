#!/bin/bash

#
# This script demonstrates building a simple feeds website from an
# EPrints repository using the _epgo_ and a few other tools developed
# by the Caltech Library.
#

function CheckSoftware () {
    NAME=$1
    PROG=$(which $NAME)
    echo -n "Checking for $NAME ... "
    if [ "$PROG" = "" ]; then
        echo "$NAME not found."
        return 1
    fi
    echo "using $PROG"
    return 0
}

function BuildIndex () {
    epgo -export $1 
}

function BuildFeeds () {
    epgo -build
}

CheckSoftware epgo && BuildIndex 10000
CheckSoftware epgo && BuildFeeds

