
# Action items

## Bugs

+ [ ] Missing committee and advisor data from thesis harvest

## Next

+ [ ] Update specific structure targeting CaltechAUTHORS EPrint repository to general purpose EPrint repository (e.g. CaltechTHESIS)
+ [ ] Depreciate _ep -select_ in favor of _dataset_ command with filter options
+ [ ] Depreciate _epgo-genpages_ in favor of _dataset_ plus _mkpage_
+ [ ] Add Authentication support for harvesting all EPrints records if the REST API is enabled
    + [ ] Use privileged account for Harvest so we can get all content
    + [ ] Export should have options to include/exclude embargoed/restricted records
+ [ ] Calculate a "changed since" return of records
    + Use fielded REST API calls to pull out the change dates with EPrint ID, then calculate the subset of records in date range requested

## Someday, Maybe

+ export should work on ID ranges as well as modified dates so we can pickup changes frequently
+ export lists of groups, funders, EPrint object types, and other fields that might be useful for filtering/sorting output
+ ORCID person outputs need to include name (e.g. could do a look up via ORCID API)
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

+ [x] Rename _epgo_ to _ep_
+ [x] When harvesting recent 1200 articles, collection.json and keys.json are being cleared (switch from dataset.Create() to dataset.Open())
+ [x] saving select lists when storage is S3 whipes out collections.keymap and keys.json
+ [x] move BuildSite() into cmds/epgo-genpages/epgo-genpages.go
+ [x] Save raw EPrint XML with harvested EPrint
+ [x] convert epgo-genpages HTML/HTML Include templates into _mkpage_ friendly templates
+ [x] Replace templated RSS2 output with output generated from epgo-genpages to _mkpage_ or other system
