
# Action items

## Bugs

+ [ ] On some repositories <note> should be suppressed in public view, in others it is public (e.g. CaltechAUTHORS Note holds copyright info)

## Next

+ [ ] Add support to create/update an full EPrint record from JSON 
+ [ ] Add support to update EPrint record's attribute paths where we have permission (e.g. add ORCID to creator objects)


## Someday, Maybe

+ [ ] Add support to show CURL version of action without running command
+ [ ] Normalize logging between ep and other harvesters
+ [ ] Add write support to ep cli so we can full circle data from other sources into EPrints
+ [ ] Add Authentication support for harvesting all EPrints records if the REST API is enabled
    + [ ] Use privileged account for Harvest so we can get all content
    + [ ] Export should have options to include/exclude embargoed/restricted records
+ [ ] Rethink how I have named the elements of the EPrint Document and Record structures. Should they more closely represent their XML structures?
    + What about versions of EPrints' data structures (e.g. how many name spaces do I need to support? http://eprints.org/ep2/data/2.0 or others?)
    + Should Record be named EPrints and Document be named EPrint?

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
+ [x] export should work on modified dates as well as ID ranges so we can pickup changes frequently
    + [x] Calculate a "changed since" return of records
        + Use fielded REST API calls to pull out the change dates with EPrint ID, then calculate the subset of records in date range requested
+ [x] Add option to export a list of keys (one per line) for exported records
    + this would let you streaming line ORCID harvests for changed records versus whole dataset collection

