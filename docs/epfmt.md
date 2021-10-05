
USAGE
=====

	epfmt [OPTIONS]

SYNOPSIS
--------


_epfmt_ is a command line program for 
pretty printing EPrint XML. It can also convert
EPrint XML to and from JSON.


DESCRIPTION
-----------


_epfmt_ parses EPrint XML (or JSON version) from
standard input and pretty prints the result to 
standard out. You can change output format XML 
and JSON by using either the '-xml' or '-json' 
option. The XML representation is based on EPrints 
3.x.  _epfmt_ does NOT interact with the EPrints API 
only the the document presented via standard
input.


OPTIONS
-------

Below are a set of options available.

```
    -e, -examples        display examples
    -generate-manpage    generate man page
    -generate-markdown   generate Markdown documentation
    -h, -help            display help
    -i, -input           input file name (read the URL connection string from the input file
    -json                output JSON version of EPrint XML
    -l, -license         display license
    -nl, -newline        if true add a trailing newline
    -o, -output          output file name
    -quiet               suppress error messages
    -v, -version         display version
    -xml                 output EPrint XML
```


EXAMPLES
--------


Pretty print EPrint XML as XML.

```
    epfmt < 123.xml
```

Pretty print from EPrint XML as JSON

```
    epfmt -json < 123.xml
```

Render EPrint JSON as EPrint XML.

```
    epfmt -xml < 123.json
```

Render EPrint simplified JSON from EPrint XML.

```
    epfmt -simplified < 123.xml
```

_epfmt_ will first parse the XML or JSON 
presented to it and pretty print the output 
in the desired format requested. If no 
format option chosen it will pretty print 
in the same format as input.


epfmt v0.1.10
