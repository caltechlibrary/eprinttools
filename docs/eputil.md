
# USAGE

	eputil [OPTIONS]

## SYNOPSIS


_eputil_ is a command line program for exploring 
EPrint REST API and EPrint XML document structure
in XML as well as JSON.


## DESCRIPTION


_eputil_ parses XML content retrieved from 
EPrints 3.x. REST API. It will render 
results in JSON or XML.  With the `-raw`
option you can get the unmodified EPrintXML from the 
REST API otherwise the XML is parsed before final 
rendering as JSON or XML. It requires a basic knowledge
of the layout of EPrint 3.x's REST API. It supports
both unauthenticated and basic authentication access
to the API. The REST API authentication mechanism 
appears indepent of the primary website authentication
setup of the installed EPrints (at least at Caltech
Library). See the examples to start exploring the API.


## OPTIONS

Below are a set of options available.

```
    -auth                       set the authentication type for access
    -e, -examples               display examples
    -generate-manpage           generate man page
    -generate-markdown          generate Markdown documentation
    -h, -help                   display help
    -json                       attempt to parse XML into generaic JSON structure
    -l, -license                display license
    -nl, -newline               if true add a trailing newline
    -o, -output                 output file name
    -pw, -password              set the password for authenticated access
    -quiet                      suppress error messages
    -raw                        get the raw EPrint REST API response
    -u, -un, -user, -username   set the username for authenticated access
    -v, -version                display version
```


## EXAMPLES


Fetch the raw unmarshaled EPrint XML via the 
EPrint REST API for id 123.

```
    eputil -raw https://example.org/rest/eprint/123.xml
```

Fetch the EPrint XML marshaled as XML using the 
EPrints REST API for id 123.

```
    eputil https://example.org/rest/eprint/123.xml 
```

Fetch the EPrint XML marshaled as JSON using the
EPrints REST API for id 123.

```
    eputil -json https://example.org/rest/eprint/123.xml
```

Get a JSON array of eprint ids from the REST API

```
    eputil -json https://example.org/rest/eprint/ 
```

Get the last modified date for id 123 from REST API

```
    eputil -raw https://example.org/rest/eprint/123/lastmod.txt 
```

If the EPrint REST API is protected by basic auth then
you can pass the username and password via the URL.
In this example the username is "user" and password is
"secret".

```
    eputil https://user:secret@example.org/rest/eprint/123.xml
```

Getting IDs doesn't typically require authentication but seeing
specific records may depending on the roles and security
setup implemented in the EPrint instance.



eputil v0.0.20
