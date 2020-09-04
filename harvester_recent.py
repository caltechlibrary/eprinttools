#!/usr/bin/env python3

#
# harvester.py is a proof of concept replacment of the `ep` command
# using Python, eputil and the eprints3x Python module. 
#

import os
import sys
import json

from eprinttools import harvest_init, harvest, harvest_keys
from eprinttools import Configuration

def usage():
    app = os.path.basename(sys.argv[0])
    print(f'''
USAGE: {app} CONGIF_JSON

{app} harvests EPrints recent record(s) and related
documents. The number of days to harvest is set in the
configuraiton file. Initally all keys are retrieved
than the last modified date is exammed via the EPrints
REST URL and records which were created or modified
in the last N days are harvested.

{app} also converts the EPrintXML into JSON which is
stored in a dataset collection and stores the related
records and EPrintXML as attachments to the JSON reccord.
If the configuration has 'include_documents' set to true
then the documents are harvested too.

  {app} config.json

This command harvests EPrint records for the last N days
described in the config.json file.

''')

if __name__ == "__main__":
    f_name = ''
    number_of_days = 0
    keys = []
    if len(sys.argv) < 2:
        usage()
        sys.exit(1)
    if len(sys.argv) >= 2:
        f_name = sys.argv[1]
    if len(sys.argv) > 2:
        for key in sys.argv[2:]:
            keys.append(key)

    if f_name == '':
        print(f'ERROR: Missing configuration filename.')
        sys.exit(1)
    if not os.path.exists(f_name):
        print(f'ERROR: Missing {f_name} file.')
        sys.exit(1)
    cfg = Configuration()
    if cfg.load_config(f_name) and cfg.required(['dataset', 'eprint_url', 'number_of_days', 'include_documents']):
        # Initialize the connection information (e.g. authentication)
        err = harvest_init(cfg.dataset, cfg.eprint_url)
        if err != '':
            print(err)
            sys.exit(1)
        if len(keys) == 0:
            keys = harvest_keys()
            if len(keys) == 0:
                print("No keys found")
                sys.exit(1)
        err = harvest(keys, include_documents = cfg.include_documents, number_of_days = cfg.number_of_days)
        if err != '':
            print(err)
            sys.exit(1)
        print('OK')
    else:
        sys.exit(1)
