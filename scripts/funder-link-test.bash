#!/bin/bash

finddir -depth 3 htdocs/funder | while read TITLE; do 
    LINK=$(urlencode $TITLE)
    echo "+ [$TITLE]($LINK)"; 
done > htdocs/funder/index.md
mkpage \
    "content=htdocs/funder/index.md" \
    templates/default/index.html \
    > htdocs/funder/index.html
