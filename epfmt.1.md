---
title: "epfmt (1) user manual"
author: "R. S. Doiel"
pubDate: 2023-01-11
---

# NAME

epfmt

# SYNOPSIS

epfmt 

# DESCRIPTION

epfmt is a command line program for pretty printing
EPrint XML. It can also convert EPrint XML to and from
JSON. By default it reads from standard input and writes to
standard out.

epfmt EPrint XML (or JSON version) from
standard input and pretty prints the result to 
standard out. You can change output format XML 
and JSON by using either the '-xml' or '-json' 
option. The XML representation is based on EPrints 
3.x.  epfmt does NOT interact with the EPrints API 
only the the document presented via standard
input.

# OPTIONS

-help
: display help

-license
: display license

-i, -input
: (string) input file name (read the URL connection string from the input file

-json
: output JSON version of EPrint XML

-nl, -newline
: if true add a trailing newline

-o, -output
: (string) output file name

-quiet
: suppress error messages

-s, -simplified
: output simplified JSON version of EPrints XML

-version
: display version

-xml
: output EPrint XML

# EXAMPLES

Pretty print EPrint XML as XML.

~~~
    epfmt < 123.xml
~~~

Pretty print from EPrint XML as JSON

~~~
    epfmt -json < 123.xml
~~~

Render EPrint JSON as EPrint XML.

~~~
    epfmt -xml < 123.json
~~~

epfmt will first parse the XML or JSON 
presented to it and pretty print the output 
in the desired format requested. If no 
format option chosen it will pretty print 
in the same format as input.

epfmt 1.2.1


