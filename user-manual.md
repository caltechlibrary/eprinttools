
# command line tools

- [eputil](eputil.1.md) is the new harvester, it talks to the REST API and can return JSON documents, XML or EPrints documents available from the REST API. It can be piped into a [dataset](https://caltechlibrary.github.io/dataset) if needed
- [epfmt](epfmt.1.md) is a command line utility to pretty print EPrints XML and convert to/from JSON
    - in the process of pretty printing it also validates the EPrints XML against the eprinttools Go package definitions
- [doi2eprintxml](doi2eprintxml.1.md) is a CaltechAUTHORS centric DOI to EPrint XML document generator 
- [ep3apid](ep3apid.1.md) is a web service that runs on localhost providing an extended EPrints 3 API via interaction with EPrints MySQL database(s)
- [ep3harvester](ep3harvester.1.md) harvests EPrints 3.x repository into a dataset collection

## Tutorials

- Converting DOIs to EPrints XML
    - [Windows 10 Workflow](windows-10-workflow.md)
    - [macOS Workflow](macos-workflow.md)
- [EPrints extended API](EPrints-extended-API.md)
- [Enabling REST API](Enabling-REST-API.md)
- [Testing setup](Test-Setup.md)
- [Composit Views](composit-views.md)
- [Generasting an EPrints user list](generating-an-eprints-user-list.md)
- [ORCID propagator](orcid-propagator.md)
- [Replacating an EPrints Repository](replicating_an_eprints_repo.md)
- [Repositories End Point](repositories-endpoint.md)
- [Repository End Point](repository-endpoint.md)
- [SQL Table Notes](sql-table-notes.md)
- [View Notes](view_notes.md)
