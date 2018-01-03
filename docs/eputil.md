
# USAGE

	eputil [OPTIONS]

## SYNOPSIS


	eputil parses XML content retrieved from disc or the EPrints API. It will 
	render JSON if the XML is valid otherwise return errors.


## OPTIONS

```
    -document, -eprints       parse an eprints (e.g. rest response) document
    -e, -examples             display examples
    -generate-markdown-docs   output documentation in Markdown
    -h, -help                 display help
    -i, -input                input file name
    -json                     attempt to parse XML into generaic JSON structure
    -l, -license              display license
    -nl, -newline             if true add a trailing newline
    -o, -output               output file name
    -paths                    get a list of doc paths (e.g. ids or sub-fields depending on the URL provided
    -quiet                    suppress error messages
    -revision, -eprint        parse a eprint (revision) document
    -url                      do an HTTP GET to fetch the XML from the URL then parse
    -v, -version              display version
```


## EXAMPLES


Fetch an EPrints document as JSON from a URL for an EPrint with an id of 123

    eputil -url https://eprints.example.org/rest/eprint/123.xml -json

Fetch an EPrints document as XML from a URL for an EPrint with an id of 123

    eputil -url https://eprints.example.org/rest/eprint/123.xml

Fetch the creators.xml as JSON for an EPrint with the id of 123.

    eputil -url https://eprints.example.org/rest/eprint/123/creators.xml -json

Parse an EPrint reversion XML document

    eputil -i revision/2.xml -eprint


eputil v0.0.10-beta7
