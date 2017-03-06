
# Action items

## Next

+ [ ] Update how BibTeX and JSON blobs are named, need to support articles, publications, thesis from two EPrints repositories
+ [ ] ORCID person outputs need to include name (e.g. could do a look up via ORCID API)
+ [ ] Remove epgo dependency on tmplfn package in favor of _mkpage_ template rendering

## Someday, Maybe

+ ORCID A-Z list
    + for each ORCID harvest as public ORCID profile write out to $ORCID_ID/orcid-profile.json
    + From orcid-profile.json rendering a Markdown document $ORCID_ID/index.md
    + Render index.html and index.include from $ORCID_ID/index.md
+ Parallelize epgo export and epgo-genpages
    + See https://gobyexample.com/worker-pools for worker pool example
    + See https://gobyexample.com/rate-limiting for rate limitting example
+ Debug Person and Group feeds with dataset integration between CaltechAUTHORS and CaltechTHESIS
+ Add feeds for Grant Numbers (use CrossREF API for naming grant sources)
+ Integrate VIAF numbers (viaf.org) and links to other numbers (e.g. ISNI)
+ Add support for the MLA citation style in the HTML markup for HTML found in browse by person/year in both the website and include files
    + Does it make more sense to include BibTeX and let BibTeX format different citation formats?


## Completed

+ [x] convert epgo-genpages HTML/HTML Include templates into _mkpage_ friendly templates
+ [x] Replace templated RSS2 output with output generated from epgo-genpages to _mkpage_ or other system
