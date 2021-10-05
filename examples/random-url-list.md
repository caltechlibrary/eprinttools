Random EPrints Sample
=====================

Problem: Select a random sample of 100 EPrint record urls from CaltechAUTHORS

Summary solution: use `eputils` to fetch a list of all eprint ids in the EPrints repository. Use `jsonrange` to take a sample of those ids
to then format the output as a list of URLs.

Detail Solution: 

1. This solution relies on Bash which is available for macOS, Linux and Windows 10 (using the Linux Sub-system)
2. Make sure you have [eprinttools](https://github.com/caltechlibrary/eprinttools/releases) version v1.0.0 or later installed
3. Make sure you have [datatools](https://github.com/caltechlibrary/datatools/releases) version v1.0.5-dev-1 or later installed.
4. Use `eputil` from eprinttools to create a JSON document of all the eprint ids in the repository
5. Use `jsonrange` from datatools to generate a random sample of those ids saving them one id per line in a file
6. Read in the list of saved ids and write out a list of URLs.

From Bash the script below prompts for the hostname of the EPrints repository as well as username and password to access the repository's REST API. Note not all EPrints accounts get access to the
REST API. Contact DLD if you need that access added for your EPrints account in Caltech Library.

The script does the following

1. Prompt for eprint hostname (e.g. https://authors.library.caltech.edu)
2. Prompt for eprint username
3. Prompt for password
4. Prompt for sample size (needs to be an integer greater than zero)

From their the script will retrieve the eprint ids, save them in a file called "eprint_ids.json", The script then reads that JSON array
and takes a sample saving the results as one ID per line in "eprint_ids.keys". Finally "eprint_ids.keys" is read in line by line
and outputs a URL pointing at the eprint id's page.

```bash
#!/bin/bash
read -p "Enter the EPrints hostname: " EP_HOST
read -p "Enter the EPrints username: " EP_USER
read -sp "Enter the EPrints password: " EP_PASSWD
read -p "Enter sample size (integer greater than zero): " SAMPLE_SIZE
eputil -auth=basic -json "https://${EP_USER}:${EP_PASSWD}@${EP_HOST}/rest/eprint/" >eprint_ids.json
cat eprint_ids.json | jsonrange -sample=${SAMPLE_SIZE} -values \
  >eprint_ids.keys
for KEY in $(cat eprint_ids.keys); do
  echo "https://${EP_HOST}/${KEY}"
done
```


