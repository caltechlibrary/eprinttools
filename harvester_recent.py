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
USAGE: {app} CONGIF_JSON [NUMBER_OF_DAYS]

{app} harvests EPrints recent record(s) and related
documents. It filters the list of all keys looking for
once with a last modified date via the REST URL.
It converts the EPrintXML into JSON which is
stored in a dataset collection and stores the related
records and EPrintXML as attachments to the JSON reccord.

  {app} config.json 5

This would harvest EPrint records last modified in the
previous five days as described in the config.json file.

''')

if __name__ == "__main__":
    f_name = ''
    c_name = ''
    url = ''
    number_of_days = 0
    keys = []
    if len(sys.argv) < 2:
        usage()
        sys.exit(1)
    if len(sys.argv) >= 2:
        f_name = sys.argv[1]
    if len(sys.argv) == 3 :
        number_of_days = int(sys.argv[2])

    if f_name == '':
        print(f'ERROR: Missing configuration filename.')
        sys.exit(1)
    if not os.path.exists(f_name):
        print(f'ERROR: Missing {f_name} file.')
        sys.exit(1)
    cfg = Configuration()
    if cfg.load_config(f_name) and cfg.required(['dataset', 'eprint_url', 'number_of_days']):
        # Initialize the connection information (e.g. authentication)
        c_name, url, number_of_days = cfg.dataset, cfg.eprint_url, cfg.number_of_days
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
        err = harvest(keys, include_documents = False, number_of_days = number_of_days) # , save_exported_keys = f'exported-recent-{repo_name}.keys')
        if err != '':
            print(err)
            sys.exit(1)
        print('OK')
    else:
        sys.exit(1)
