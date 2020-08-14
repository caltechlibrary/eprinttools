#!/usr/bin/env python3

import sys
import os

from eprints3x import s3_publish, Logger

log = Logger(os.getpid())


#
# Main processing
#
if __name__ == '__main__':
    f_name = 'config.json'
    args = []
    if len(sys.argv) > 1:
        f_name = sys.argv[1]
    if len(sys.argv) > 2:
        args = sys.argv[2:]
    s3_publish(f_name, args)

