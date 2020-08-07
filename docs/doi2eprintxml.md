
# USAGE

	doi2eprintxml [OPTIONS] DOI

## DESCRIPTION


doi2eprintxml is a Caltech Library centric application that
takes one or more DOI, queries the CrossRef API
and if that fails the DataCite API and returns an
EPrints XML document suitable for import into
EPrints. The DOI can be in either their canonical
form or URL form (e.g. "10.1021/acsami.7b15651" or
"https://doi.org/10.1021/acsami.7b15651").



## OPTIONS

Below are a set of options available.

```
    -c, -crossref        only search CrossRef API for DOI records
    -clsrules            Apply Caltech Library Specific Rules to EPrintXML output
    -d, -datacite        only search DataCite API for DOI records
    -eprints-url         Sets the EPRints API URL
    -generate-manpage    generate man page
    -generate-markdown   generate Markdown documentation
    -h, -help            display help
    -i, -input           set input filename
    -json                output EPrint structure as JSON
    -l, -license         display license
    -m, -mailto          set the mailto value for CrossRef API access
    -quiet               set quiet output
    -v, -version         display app version
```


## EXAMPLES


Example generating an EPrintsXML for one DOI

	doi2eprintxml "10.1021/acsami.7b15651" > article.xml

Example generating an EPrintsXML for two DOI

	doi2eprintxml "10.1021/acsami.7b15651" "10.1093/mnras/stu2495" > articles.xml

Example processing a list of DOIs in a text file into
an XML document called "import-articles.xml".

	doi2eprintxml -i doi-list.txt -o import-articles.xml


doi2eprintxml v0.0.58
