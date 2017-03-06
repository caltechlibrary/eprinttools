
# Action items

## Next

+ [ ] Update how BibTeX and JSON blobs are written, I need to support articles, publications, thesis from two EPrints repositories with same utility
    + [ ] provide a mechanism to create groupings of EPrint records to handle Eprint Object types (e.g. articles vs. thesis), groups, orcids, etc.
    + [ ] move BuildSite() into cmds/epgo-genpages/epgo-genpages.go
+ [ ] ORCID person outputs need to include name (e.g. could do a look up via ORCID API)
+ [ ] Remove epgo dependency on tmplfn package in favor of _mkpage_ template rendering

## Someday, Maybe

+ export single EPrint record to dataset so we can do fast additions on breaking publications
+ export lists of groups, funders, EPrint object types, and other fields that might be useful for filtering/sorting output
+ ORCID A-Z list
    + for each ORCID harvest as public ORCID profile write out to $ORCID_ID/orcid-profile.json
    + From orcid-profile.json rendering a Markdown document $ORCID_ID/index.md
    + Render index.html and index.include from $ORCID_ID/index.md
+ Find epgo-genpage bottleneck and improve performance
+ Debug Person and Group feeds with dataset integration between CaltechAUTHORS and CaltechTHESIS
+ Add feeds for Grant Numbers (use CrossREF API for naming grant sources)
+ Integrate VIAF numbers (viaf.org) and links to other numbers (e.g. ISNI)
+ Add support for the MLA citation style in the HTML markup for HTML found in browse by person/year in both the website and include files
    + Does it make more sense to include BibTeX and let BibTeX format different citation formats?


## Completed

+ [x] convert epgo-genpages HTML/HTML Include templates into _mkpage_ friendly templates
+ [x] Replace templated RSS2 output with output generated from epgo-genpages to _mkpage_ or other system
