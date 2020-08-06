
# USAGE

	apiexplorer [OPTIONS]

## SYNOPSIS


_apiexplorer_ is a demo of rc package for accessing REST API


## DESCRIPTION


_apiexploere_ is a demo program exercising the rc Golang package.


## OPTIONS

Below are a set of options available.

```
    -as-json            Convert XML to JSON before output
    -auth               set authorization type (e.g. oauth, shib)
    -generate-manpage   generate man page
    -generate-markdown  generate markdown documentation
    -h, -help           display help
    -l, -license        display license
    -method             set the http method to use for request
    -o, -output         output filename
    -payload            A JSON structure holding the payload data
    -pw, -password      set user secret to use for authentication
    -un, -username      set username for authentication
    -v, -version        display version
```


## EXAMPLES


_apiexplorer_ accessing a public DataCite API for DOI
"10.22002/D1.924".

```
    apiexplorer -method GET \
    "https://api.datacite.org/works/10.22002/D1.924?email=jdoe@example.edu"
```

	

apiexplorer v0.0.1
