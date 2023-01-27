---
title: "doi2eprintxml"
author: "R. S. Doiel"
pubDate: 2023-01-11
---

# NAME

doi2eprintxml

# SYNOPSIS

doi2eprintxml [OPTIONS] DOI

# DESCRIPTION

doi2eprintxml is a Caltech Library centric application that
takes one or more DOI, queries the CrossRef API
and if that fails the DataCite API and returns an
EPrints XML document suitable for import into
EPrints. The DOI can be in either their canonical
form or URL form (e.g. "10.1021/acsami.7b15651" or
"https://doi.org/10.1021/acsami.7b15651").

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-D
: attempt to download the digital object if object URL provided

-c
: only search CrossRef API for DOI records

-clsrules
: Apply current Caltech Library Specific Rules to EPrintXML output (default true)

-crossref
: only search CrossRef API for DOI records

-d
: only search DataCite API for DOI records

-datacite
: only search DataCite API for DOI records

-dot-initials
: Add period to initials in given name

-download
: attempt to download the digital object if object URL provided

-eprints-url string
: Sets the EPRints API URL

-i, -input
: (string) set input filename

-json
: output EPrint structure as JSON

-m 
: (string) set the mailto value for CrossRef API access (default "helpdesk@library.caltech.edu")

-mailto
: (string) set the mailto value for CrossRef API access (default "helpdesk@library.caltech.edu")

-normalize-publisher
: Use normalize publisher rule

-normalize-related-url
: Use normlize related url rule

-normlize-publication
: Use normalize publication rule

-o, -output
: (string) set output filename

-quiet
: set quiet output

-simple
: output EPrint structure as Simplified JSON

-trim-creators
: Use trim creators list rule

-trim-number
: Use trim number rule

-trim-series
: Use trim series rule

-trim-title
: Use trim title rule

-trim-volume
: Use trim volume rule

# EXAMPLES

Example generating an EPrintsXML for one DOI

~~~
	doi2eprintxml "10.1021/acsami.7b15651" > article.xml
~~~

Example generating an EPrintsXML for two DOI

~~~
	doi2eprintxml "10.1021/acsami.7b15651" "10.1093/mnras/stu2495" > articles.xml
~~~

Example processing a list of DOIs in a text file into
an XML document called "import-articles.xml".

~~~
	doi2eprintxml -i doi-list.txt -o import-articles.xml
~~~

doi2eprintxml 1.2.1


