
# USAGE

	eprintxml2json [OPTIONS] [input filename]

## SYNOPSIS

_eprintxml2json_ convers EPrintXML documents to JSON

## DESCRIPTION

_eprintxml2json_ converts EPrintXML like that
retrieved from the EPrint 3.x REST API to JSON. If no filename
is provided on the command line then standard input is used
to read the EPrint XML. If the EPrint XML isn't understood
then an error message will be written and an exit code of 1
used to close the process otherwise the process will render
JSON to standard out.


## OPTIONS

Below are a set of options available.

```
    -e, -examples        display examples
    -generate-manpage    generate man page
    -generate-markdown   generate Markdown documentation
    -h, -help            display help
    -l, -license         display license
    -nl, -newline        if true add a trailing newline
    -o, -output          output file name
    -p, -pretty          pretty print output
    -quiet               suppress error messages
    -v, -version         display version
```


## EXAMPLES

Converting a document, eprints-dump.xml, to JSON.

```
    eprintxml2json eprints-dump.xml
```

Or

```
    cat eprints-dump.xml | eprintxml2json 
```



eprintxml2json v0.0.58
