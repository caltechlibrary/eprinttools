#!/usr/bin/env python3

import sys
import os

from eprints3x import s3_publish, Logger
from eprintviews import Configuration

log = Logger(os.getpid())


#
# Main processing
#
if __name__ == '__main__':
    f_name = ''
    args = []
    if len(sys.argv) > 1:
        f_name = sys.argv[1]
    if len(sys.argv) > 2:
        args = sys.argv[2:]
    if f_name == '':
        print(f'Missing configuration filename.')
        sys.exit(1)
    if not os.path.exists(f_name):
        print(f'Missing {f_name} file.')
        sys.exit(1)
    cfg = Configuration()
    if cfg.load_config(f_name) and cfg.required(['htdocs', 'bucket']):
        htdocs, bucket = cfg.htdocs, cfg.bucket
        s3_publish(htdocs, bucket, args)
    else:
        sys.exit(1)
