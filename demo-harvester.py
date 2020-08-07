#!/usr/bin/env python3

#
# harvester.py is a proof of concept replacment of the `ep` command
# using Python, eputil and the eprints3x Python module. 
#

import os
import sys

from eprints3x import harvest_init, harvest, harvest_keys

def usage():
    app = os.path.basename(sys.argv[0])
    print(f'''
USAGE: {app} DATASET_COLLECTION EPRINT_URL [KEY [KEY] ...]

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
    if len(sys.argv) < 3:
        usage()
        sys.exit(1)
    if len(sys.argv) >= 3:
        c_name, url = sys.argv[1], sys.argv[2]
    if len(sys.argv) > 3:
        for key in sys.argv[3:]:
            keys.append(key)

    # Initialize the connection information (e.g. authentication)
    err = harvest_init(c_name, url)
    if err != '':
        print(err)
        sys.exit(1)
    if len(keys) == 0:
        keys = harvest_keys()
        if len(keys) == 0:
            print("No keys found")
            sys.exit(1)
    repo_name, _ = os.path.splitext(c_name)
    err = harvest(keys, include_documents = True, save_exported_keys = f'exported-{repo_name}.keys')
    if err != '':
        print(err)
        sys.exit(1)
    print('OK')
