#!/usr/bin/env python3

#
# epharvest.py is a prototoype using eputil command to replace
# ep command for harvesting EPrintXML, document objects and 
# storing them in a dataset collection.
#

import os
import sys

from eprints3x import init, harvest_keys, harvest_record


def usage():
    app = os.path.basename(sys.argv[0])
    print(f'''
USAGE: {app} DATASET_COLLECTION EPRINT_URL [KEY]

{app} harvests EPrints record(s) and related
documents. It converts the EPrintXML into JSON which is
stored in a dataset collection and stores the related
records and EPrintXML as attachments to the JSON reccord.

  {app} repository.ds \\
     'https://automation:SECRET@repository.example.edu' \\
     123

This would harvest EPrint record 123 from repository.example.edu
EPrints server.

''')

if __name__ == "__main__":
    c_name = ''
    url = ''
    keys = []
    if len(sys.argv) == 3:
        c_name, url = sys.argv[1], sys.argv[2]
    elif len(sys.argv) == 4:
        c_name, url, key = sys.argv[1], sys.argv[2], sys.argv[3]
        keys = [ key ]
    else:
        usage()
        sys.exit(1)

    # Initialize the connection information (e.g. authentication)
    base_url, err = init(c_name, url)
    if err != '':
        print(err)
        sys.exit(1)
    if len(keys) == 0:
        keys = harvest_keys(base_url)
        if len(keys) == 0:
            print("No keys found")
            sys.exit(1)
    for key in keys:
        err = harvest_record(base_url, key)
        if err != '':
            print(err)
            sys.exit(1)
    print('OK')
