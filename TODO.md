
# Action items

## Next

+ Figure out why select lists are empty!
+ Add option to epgo to regenerate select lists from collection
+ epgo-genpages - rethink how generate markdown, BibTex, JSON, RSS2 only, this would drop the need for epgo specific templates
+ epgo-indexer - should operete by ingesting Markdown an selected JSON fields and pointing at HTML
    + could we just ingest Markdown and point at HTML? This would let the indexer be generic and part of _mkpage_
+ epgo-sitemapper needs to be depreciated in favor of by _mkpage_ sitemapper
+ Add ORCID API harvest for person biographies and profiles, populating the person landing page with this content
    + if we harvest profiles form the new ORCID API as orcid-profile.json I can use csvcols to extract and assemble addition fields for profile landing page
+ ORCID A-Z list
    + for each ORCID harvest as public ORCID profile write out to $ORCID_ID/orcid-profile.json
    + From orcid-profile.json rendering a Markdown document $ORCID_ID/index.md
    + Render index.html and index.include from $ORCID_ID/index.md

## Someday, Maybe

+ Parrallelize epgo export and epgo-genpages
    + See https://gobyexample.com/worker-pools for worker pool example
    + See https://gobyexample.com/rate-limiting for rate limitting example
+ Debug Person and Group feeds with dataset integration between CaltechTHESIS, CaltechARCHIVES and CaltechAUTHORS
+ Add feeds for Grant Numbers
+ Integrate ISNI content
+ Integrate VIAF numbers (viaf.org) and links to other numbers (e.g. ISNI)
+ Add support for the MLA citation style in the HTML markup for HTML found in browse by person/year in both the website and include files
    + Does it make more sense to include BibTeX and let BibTeX format different citation formats?


