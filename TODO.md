
# Action items

## Next

+ convert epgo-genpages HTML/HTML Include templates into _mkpage_ friendly templates
+ Add ORCID API harvest for person biographies and profiles, populating the person landing page with this content
    + if we harvest profiles form the new ORCID API as orcid-profile.json I can use csvcols to extract and assemble addition fields for profile landing page

## Someday, Maybe

+ Replace templated RSS2 output with output generated form rss2 package
    + this would allow removing the dependency on tmplfn and text/template
+ ORCID A-Z list
    + for each ORCID harvest as public ORCID profile write out to $ORCID_ID/orcid-profile.json
    + From orcid-profile.json rendering a Markdown document $ORCID_ID/index.md
    + Render index.html and index.include from $ORCID_ID/index.md
+ Parallelize epgo export and epgo-genpages
    + See https://gobyexample.com/worker-pools for worker pool example
    + See https://gobyexample.com/rate-limiting for rate limitting example
+ Debug Person and Group feeds with dataset integration between CaltechTHESIS, CaltechARCHIVES and CaltechAUTHORS
+ Add feeds for Grant Numbers
+ Integrate ISNI content
+ Integrate VIAF numbers (viaf.org) and links to other numbers (e.g. ISNI)
+ Add support for the MLA citation style in the HTML markup for HTML found in browse by person/year in both the website and include files
    + Does it make more sense to include BibTeX and let BibTeX format different citation formats?


