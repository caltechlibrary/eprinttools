#!/usr/bin/env python3

import csv
import json
import os
from subprocess import Popen, PIPE, run
import sys

from .logger import Logger

#
# publisher.py is a script to publish an htdocs directory 
# to an S3 bucket. It use specific to Linux/Mac OS X/Unix
# flavored environments.
#
log = Logger(os.getpid())

# Mime type map by file extensions
type_map = {
        ".bib": "text/plain",
        ".css": "text/css",
        ".csv": "text/csv",
        ".gif": "image/gif",
        ".gz": "application/gzip",
        ".html": "text/html",
        ".ico": "image/x-icon",
        ".include": "text/plain",
        ".js": "application/javascript",
        ".json": "application/json",
        ".keys": "text/plain",
        ".md": "text/markdown",
        ".png": "image/png",
        ".rss": "application/rss+xml",
        ".svg": "image/svg+xml",
        ".txt": "text/plain",
        ".zip": "application/zip"
    }

# Encoding type by extensions
encoding_map = {
        ".bib": "UTF-8",
        ".css": "UTF-8",
        ".csv": "UTF-8",
        ".html": "UTF-8",
        ".include": "UTF-8",
        ".js": "UTF-8",
        ".json": "UTF-8",
        ".keys": "UTF-8",
        ".md": "UTF-8",
        ".rss": "UTF-8",
        ".svg": "UTF-8",
        ".txt": "UTF-8",
        ".zip": "UTF-8"
    }




def s3_sync_by_ext(root, bucket_name, ext = "",  content_type = "", content_encoding = "", delete_removed = True, dryrun = False, quiet = True):
    cmd = [
            'aws', 
            's3',
            'sync',
            root, 
            bucket_name,
            '--acl', 'public-read'
            ]
    if quiet == True:
        cmd.append('--quiet')
    if ext != "":
        cmd.append('--exclude')
        cmd.append('*')
        cmd.append('--include')
        cmd.append(f'*{ext}')
    if content_type != '':
        cmd.append('--content-type')
        cmd.append(content_type)
    if content_encoding != '':
        cmd.append('--content-encoding')
        cmd.append(content_encoding)
    if delete_removed:
        cmd.append('--delete')
    if dryrun:
        cmd.append('--dryrun')
    log.print(f'{" ".join(cmd)}')
    with Popen(cmd, stdout=PIPE) as proc:
        for line in proc.stdout:
            log.print(line.strip().decode('utf-8'))
        log.print(f'Completed: {" ".join(cmd)}');


def s3_publish(htdocs, bucket, args):
    if not os.path.exists(htdocs):
        log.fatal(f'''Cannot find the htdocs directory''')
        sys.exit(1)
    if bucket == '':
        log.fatal(f'''Can't find bucket in {bucket} configuration file''')
    if bucket == '' or htdocs == '':
        log.fatal('publisher.py is not configured, check your configuration')

    # Now setup default actions if none provided
    available_options = [ ".txt", ".css", ".csv", ".gif", ".png", ".ico", ".svg", ".js", ".keys", ".csv", ".bib", ".rss", ".keys", ".md", ".json", ".include", ".html", ".include", ".zip" ]
    
    if len(args) == 0:
        args = available_options
    else:
        # Validate args
        for arg in args:
            if not arg in available_options:
                log.fatal(f"Don't know how to {arg}")
    log.print(f"Publishing {', '.join(args)} to {bucket}")
    if len(args) > 0:
        # Make sure we have an htdocs directory
        if os.path.exists(htdocs):
            log.print(f"Syncing {htdocs} by extensions {', '.join(args)} to {bucket}")
        else:
            log.print(f"Can't find htdocs {htdocs}")
            sys.exit(1)
        for ext in args: 
            log.print(f'Processing extension "{ext}"')
            content_type = 'application/octet-stream'
            content_encoding = ''
            if ext in type_map:
                content_type = type_map[ext]
            if ext in encoding_map:
                content_encoding = encoding_map[ext]
            s3_sync_by_ext(htdocs, bucket, 
                    ext, content_type, content_encoding, delete_removed = True,
                    dryrun = False)
        log.print(f"Syncing htdocs with {bucket}, completed")
        log.print("All Done!")
    
