
# Action items

## Bugs

## Next

+ [ ] Add Authentication support for harvesting all EPrints records if the REST API is enabled
    + [ ] Use privileged account for Harvest so we can get all content
    + [ ] Export should have options to include/exclude embargoed/restricted records
+ [ ] export should work on ID ranges as well as modified dates so we can pickup changes frequently
    + [ ] Calculate a "changed since" return of records
        + Use fielded REST API calls to pull out the change dates with EPrint ID, then calculate the subset of records in date range requested

## Someday, Maybe


## Completed

+ [x] Depreciate _ep -select_ in favor of _dataset_ command with filter options
+ [x] Depreciate _epgo-genpages_ in favor of _dataset_ plus _mkpage_
+ [x] Update specific structure targeting CaltechAUTHORS EPrint repository to general purpose EPrint repository (e.g. CaltechTHESIS)
+ [x] Missing committee and advisor data from thesis harvest
+ [x] Rename _epgo_ to _ep_
+ [x] When harvesting recent 1200 articles, collection.json and keys.json are being cleared (switch from dataset.Create() to dataset.Open())
+ [x] saving select lists when storage is S3 whipes out collections.keymap and keys.json
+ [x] move BuildSite() into cmds/epgo-genpages/epgo-genpages.go
+ [x] Save raw EPrint XML with harvested EPrint
+ [x] convert epgo-genpages HTML/HTML Include templates into _mkpage_ friendly templates
+ [x] Replace templated RSS2 output with output generated from epgo-genpages to _mkpage_ or other system
