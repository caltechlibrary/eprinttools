
# Action items

## Bugs

+ [x] Subjects and Keywords need to normalize before generating the index.json file, they are showing up as "true" in processed by Pandoc.
    + I have normalized these as subject_list and keyword_list so we retain the EPrintsXML mapping but have useful data organization for Pandoc.
+ [x] We need to have an optional base_url in the templates and mk_website.py's assemble so that we can handle the Archives usecase was well as 4th level domain website use case
+ [ ] Generalize indexer.py and static/js/searchbox.js to support all repositories scheme.py
+ [ ] Make sure rendered pages point at the S3 copy of the files rather than the original EPrint location
+ [ ] Site must meet ADA sec. 503 accessibility requirements, fix templates and CSS accordingly

## Next

+ [ ] Release v0.1.9 to pypi and GitHub
+ [ ] Add alternative sort options for aggregations via frames
    + Currently provide sort by publication date and title
    + E.g. Descending publication date sort by sort_name
+ [ ] Come up with appropriate scheme.json structure that works across repositories, then move code into eprinttools/views/scheme.py
+ [x] Build out eprints3x Python 3 module for interactiing with EPrints REST API
    + [x] eprints
    + [x] s3_publisher (we need this when harvesting document objects directly to S3)
    + [x] logger
    + [x] mysql_access (optional MySQL fast access to recently modified keys)
+ [x] Build out eprintviews module
    + [x] Configure
    + [x] Views
    + [x] Users
    + [x] Subjects
    + [x] Aggregator
+ [x] Build out general purpose harvester
+ [x] Build out general purpose genviews
+ [x] Build out general purpose indexer
+ [x] Build out general purpose mk-website
+ [x] Build out general purpose publisher
+ [x] Build out cloud front invalidator example
+ [ ] Rather than harvest digital object files into dataset collection, put them directly in S3 in their appropriate path
+ [ ] Identifiy, implement and integrate a stats page via the analytics avaialble for S3 based projects (needs to provide at least the numbers that IRStats does)
+ [ ] Integrate Search based on views and landing page's scheme.json (e.g. Lunr.js then Elasticsearch)
+ [ ] Integrate Universal Viewer into site replication
+ [ ] List all attached files not just the primary PDF
+ [ ] Fleshout scheme.json for indexing with Elasticsearch as well as Lunr.js
+ [ ] review index.json to be a more generalize metadata structure for rendering index.md landing pages, right now it represents EPrint XML in JSON
+ [ ] Add support for *.include, *.json, *.rss and *.bib (BibTeX) for all lists and search results
+ [ ] Add pass through support for OAI-PMH (e.g. view redirects)
+ [x] Come up with appropriate set of public URLs, where necessary rename EPrints repositories or redirect from homepages and header redirects to public view.
+ [x] Build out on demand refresh process for updating static sites
+ [ ] Document everything

## Someday, Maybe

+ [ ] Add support for a repository of repositories service as a dark view on datawork behind single signon for library staff (we probably can authenticate against Google Apps via OAuth2 since that is limited to library staff)
+ [ ] Prototype a search page that does the "Advanced" search by assembling the query terms into the Elasticsearch query language
+ [ ] Build out cross repository searching in the search page(s)
+ [ ] Build out a mechanism to integrate Tom's reporting 
+ [ ] Build out additional aggregation and views that will assist Librarians in assembling materials for Faculty/Researchers like what have they edited, reviewed, etc.
+ [ ] Issues #13, #14 at GitHub are bugs and enhancements.

## Completed

+ [x] DR-198, Issue # 28, Need to use the EPrints document version directory to set the patch level of an attachment, currently all attachments are stored as v0.1.1
+ [x] DR-198, replace ep with python module.
+ [x] Normalize logging between eputil and ep (removed ep)
+ [x] DR-118, create two synthetic fields in our JSON representation, primary_object should point at the primary resource, e.g. the PDF of an article, and the second synethetic field would be related_objects which would be an array pointing at things like supplimental materials. This will make it easier to create both a landing page as well as point directly at items when rendering formats like BibTeX.
+ [x] Jira issues DR-45 and DR-40 
+ [x] ep tool is exporting whole URL as key rather than the eprint id number
+ [x] Implement a DOI to EPrintXML cli (e.g. api.crossref.org, api.datacite.org)
+ [x] failed to export for status 'inbox' and 'deletion' should be a warning
+ [x] Upgrade to dataset v0.0.14-dev or better
+ [x] Add support to create/update an full EPrint record from JSON  (
+ [x] Rethink how I have named the elements of the EPrint Document and Record structures. Should they more closely represent their XML structures?
    + What about versions of EPrints' data structures (e.g. how many name spaces do I need to support? http://eprints.org/ep2/data/2.0 or others?)
    + Should Record be named EPrints and Document be named EPrint?
+ [x] Add support to show CURL version of action without running command
    + done in new "eprint" cli that will replace "ep"
+ [x] Add support to update EPrint record's attribute paths where we have permission (e.g. add ORCID to creator objects)
    + added via the new "eprint" cli
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
+ [x] On some repositories <note> should be suppressed in public view, in others it is public (e.g. CaltechAUTHORS Note holds copyright info)
    + If you access the repository without authenticating (i.e. the public view of the REST API) the note field is suppressed by EPrints.

