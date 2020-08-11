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
USAGE: {app} DATASET_COLLECTION EPRINT_URL NUMBER_OF_DAYS

{app} harvests EPrints recenlty modified records for the
last N days and related documents. It converts the EPrintXML 
into JSON which is stored in a dataset collection and 
stores the related records and EPrintXML as attachments 
to the JSON reccord.

  {app} repository.ds \\
     'https://USERNAME:SECRET@repository.example.edu' \\
     -14

This would harvest EPrint records modified in the last 
fourteen days from repository.example.edu
EPrints server.

''')

if __name__ == "__main__":
    c_name = ''
    url = ''
    number_of_days = 0
    keys = []
    print(f'DEBUG argc = {len(sys.argv)}')
    if len(sys.argv) < 4:
        usage()
        sys.exit(1)
    if len(sys.argv) == 4:
        c_name, url, number_of_days = sys.argv[1], sys.argv[2], sys.argv[3]
        number_of_days = int(number_of_days)
    if len(sys.argv) > 4:
        usage()
        sys.exit(1)
    if (c_name == '') or (url == '') or (number_of_days == 0):
        usage()
        sys.exit(1)

    # Initialize the connection information (e.g. authentication)
    err = harvest_init(c_name, url)
    if err != '':
        print(err)
        sys.exit(1)
    keys = harvest_keys()
    if len(keys) == 0:
        print(f"No keys found")
        sys.exit(1)
    repo_name, _ = os.path.splitext(c_name)
    number_of_days = int(number_of_days)
    if number_of_days > 0:
        number_of_days = number_of_days * -1
    err = harvest(keys, include_documents = True, save_exported_keys = f'exported-recent-{repo_name}.keys', number_of_days = number_of_days)
    if err != '':
        print(err)
        sys.exit(1)
    print('OK')
