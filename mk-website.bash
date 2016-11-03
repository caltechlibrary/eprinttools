#!/bin/bash
#

function makePage () {
    title=$1
    page=$2
    nav=$3
    html_page=$4
    echo "Generating $html_page"
    mkpage \
        "title=text:$title" \
        "content=$page" \
        "nav=$nav" \
        page.tmpl > $html_page
}


# index.html
makePage "epgo" README.md nav.md index.html

# install.html
makePage "epgo" INSTALL.md nav.md install.html

# license.html
makePage "epgo" "markdown:$(cat LICENSE)" nav.md license.html

# notes.html
makePage "epgo" NOTES.md nav.md notes.html

# todo.html
makePage "epgo" TODO.md nav.md todo.html

# Add the files to git as needed
git add index.html install.html license.html notes.html todo.html
