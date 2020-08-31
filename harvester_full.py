#!/usr/bin/env python3

#
# harvester.py is a proof of concept replacment of the `ep` command
# using Python, eputil and the eprints3x Python module. 
#

import os
import sys
import json

from eprinttools.eprints3x import harvest_init, harvest, harvest_keys
from eprinttools.eprintviews import Configuration

def usage():
    app = os.path.basename(sys.argv[0])
    print(f'''
USAGE: {app} CONGIF_JSON [KEY [KEY] ...]

{app} harvests EPrints record(s) and related
documents. It converts the EPrintXML into JSON which is
stored in a dataset collection and stores the related
records and EPrintXML as attachments to the JSON reccord.

  {app} config.json \\
     123

This would harvest EPrint record 123 from repository
described in the config.json file.

''')

if __name__ == "__main__":
    f_name = ''
    c_name = ''
    url = ''
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
    if cfg.load_config(f_name) and cfg.required(['dataset', 'eprint_url']):
        # Initialize the connection information (e.g. authentication)
        c_name, url = cfg.dataset, cfg.eprint_url
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
        err = harvest(keys, include_documents = False) #, save_exported_keys = f'exported-{repo_name}.keys')
        if err != '':
            print(err)
            sys.exit(1)
        print('OK')
    else:
        sys.exit(1)
