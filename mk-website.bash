#!/bin/bash
#

function makePage() {
	title=$1
	page=$2
	nav=$3
	html_page=$4
	echo "Generating $html_page"
	mkpage \
		"title=text:$title" \
		"content=$page" \
		"nav=$nav" \
		page.tmpl >"$html_page"
}

# index.html
makePage "eprinttools" README.md nav.md index.html

# install.html
makePage "eprinttools" INSTALL.md nav.md install.html

# license.html
makePage "eprinttools" "markdown:$(cat LICENSE)" nav.md license.html

# Add the files to git as needed
git add index.html install.html license.html

# Loop through commands docs
for FNAME in index ep eputil; do
	makePage "eprinttools" docs/$FNAME.md docs/nav.md docs/$FNAME.html
	git add docs/$FNAME.md docs/$FNAME.html
done
