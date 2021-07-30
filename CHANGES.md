
Release 1.0.1:

- Depreciating the `-clsrules=true` option. when using `doi2eprintxml`. Before release 1.0.1 the DOI was included by default in the related url field. This was an artifact of Caltech Library support DOI before EPrints did. This practice ended in 2021. Additional haviors like limiting the number of authors ended in 2020.  To reflect these changes in practive the `-clsrules` option now defaults to "false" and these rules are not applied unless you explicitly add `-clsrules` on the command line.
    - See clsrules/README.md for details
- Added in release 1.0.1 is an option to map ISSN to publisher names overriding what is provided by CrossRef and DataCite. To include include this option choose `-issn-to-publisher` option.  This option will do two things, it will create an `issn-to-publisher.csv` if none existed and if one does exist it will read it in using that mapping. In this way we can control how publishers are named.
