
# USAGE

    epgo-indexpages [OPTIONS] [BLEVE_INDEX_NAME]

## SYNOPSIS

epgo-indexpages is a command line utility to indexes content in the htdocs directory.
It produces a Bleve search index used by servepages web service.
Configuration is done through environmental variables.

## CONFIGURATION

epgo-indexpages relies on the following environment variables for
configuration when overriding the defaults:

+ EPGO_HTDOCS This should be the path to the directory tree
              containings the content (e.g. JSON files) to be index.
              This is generally populated with the caitpage command.
            Defaults to ./htdocs.

+ EPGO_BLEVE This is is the directory that will contain all the Bleve
             indexes. Defaults to ./htdocs.bleve

## OPTIONS

```
    -batch    Set the batch index size
    -bleve    a colon delimited list of Bleve index db names
    -h    display help
    -help    display help
    -htdocs    The document root for the website
    -l    display license
    -license    display license
    -r    Replace the index if it exists
    -repository-path    Path of rendered repository content
    -v    display version
    -version    display version
```


