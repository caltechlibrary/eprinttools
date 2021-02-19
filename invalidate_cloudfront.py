#!/usr/bin/env python3

import csv
import json
import os
import sys
from subprocess import Popen, PIPE, run
import sys

from eprinttools import Logger, Configuration

#
# This script invalidates cloud front CDN
#
log = Logger(os.getpid())

def invalidate_cloudfront(distribution_id):
    # aws cloudfront create-invalidation --distribution-id distribution_ID --paths "/*"
    cmd = [
        'aws',
        'cloudfront',
        'create-invalidation',
        '--distribution-id',
        distribution_id,
        '--paths',
        '/*'
    ]
    log.print(f'{" ".join(cmd)}')
    with Popen(cmd, stdout=PIPE) as proc:
        for line in proc.stdout:
            log.print(line.strip().decode('utf-8'))
        log.print(f'Completed: {" ".join(cmd)}');

#
# Main processing
#

# Make sure we have a configuration.
if __name__ == "__main__":
    app_name = os.path.basename(sys.argv[0])
    f_name = ''
    if len(sys.argv) > 1:
        f_name = sys.argv[1]
    cfg = Configuration()
    if cfg.load_config(f_name) and cfg.required(['distribution_id']):
        distribution_id = cfg.distribution_id
        log.print(f"Invalidating {distribution_id} in Cloud Front")
        invalidate_cloudfront(distribution_id)
        log.print("All Done!")
    else: 
        sys.exit(1)
