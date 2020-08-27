#!/usr/bin/env python3

#
# harvester.py is a proof of concept replacment of the `ep` command
# using Python, eputil and the eprints3x Python module. 
#

import os
import sys
import json

from eprints3x import harvest_init, harvest, harvest_keys

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
    with open(f_name) as f:
        src = f.read()
        cfg = json.loads(src)
        if 'dataset' in cfg:
            c_name = cfg['dataset']
        if 'eprint_url' in cfg:
            url = cfg['eprint_url']
    if url == '':
        print(f'ERROR: missing eprint_url in {f_name}')
        sys.exit(1)
    if c_name == '':
        print(f'ERROR: missing collection name in {f_name}')
        sys.exit(1)

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
    err = harvest(keys, include_documents = False) #, save_exported_keys = f'exported-{repo_name}.keys')
    if err != '':
        print(err)
        sys.exit(1)
    print('OK')
