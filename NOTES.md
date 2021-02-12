
# Overview

Caltech Library uses EPrints heavily for several different repositories 
(e.g. CaltechAUTHORS, CaltechTHESIS). We wished to use the metadata from 
our repositories in ways that were not convient inside EPrints as well as 
aggregate our metadata with other systems.  This lead us to create the 
_eputil_ and _ep_ command line programs making it easy get JSON 
representations of EPrint records output via EPrints' REST API.

## Reference materials

+ [Web API Docs](http://wiki.eprints.org/w/API:EPrints/Apache/CRUD)
    + includes example curl interactions
+ [Golang Template Tutorial](https://elithrar.github.io/article/approximating-html-template-inheritance/)
    + page.html and page.include are the two templates needed to produce the publications feed pages



