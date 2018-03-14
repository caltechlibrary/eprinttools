#!/usr/bin/env python3
from distutils.core import setup
from site import getsitepackages

import sys
import os
import shutil
import json

readme_md = "README.md"
readme_txt = "README.txt"

def read(fname):
    with open(fname, mode = "r", encoding = "utf-8") as f:
        src = f.read()
    return src

codemeta_json = "codemeta.json"
# If we're in the main Git repository then pickup the most recent codemeta.json and README files.
if os.path.exists(os.path.join("..", codemeta_json)):
    shutil.copyfile(os.path.join("..", codemeta_json),  codemeta_json)
if os.path.exists(os.path.join("..", readme_md)):
    shutil.copyfile(os.path.join("..", readme_md),  readme_txt)

# Let's pickup as much metadata as we need from codemeta.json
with open(codemeta_json, mode = "r", encoding = "utf-8") as f:
    src = f.read()
    meta = json.loads(src)

# Let's make our symvar string
version = meta["version"]

# Now we need to pull and format our author, author_email strings.
author = ""
author_email = ""
for obj in meta["author"]:
    given = obj["givenName"]
    family = obj["familyName"]
    email = obj["email"]
    if len(author) == 0:
        author = f"{given} {family}"
    else:
        author = author + f", {given} {family}"
    if len(author_email) == 0:
        author_email = f"{email}"
    else:
        author_email = author_email + f", {email}"

# Setup for our Go based shared library as a "data_file" since Python doesn't grok Go.
platform = os.uname().sysname
shared_library_name = "libeprint3.so"
OS_Classifier = "Operating System :: POSIX :: Linux"
if platform.startswith("Darwin"):
    shared_library_name = "libeprints3.dylib"
    platform = "Mac OS X"
    OS_Classifier = "Operating System :: MacOS :: MacOS X"
elif platform.startswith("Win"):
    shared_library_name = "libeprints3.dll"
    platform = "Windows"
    OS_Classifier = "Operating System :: Microsoft :: Windows :: Windows 10"

site_package_location = os.path.join(getsitepackages()[0], "eprinttools")

# Now that we know everything configure out setup
setup(name = "eprinttools",
    version = version,
    description = "A python module for interacting with EPrints 3 REST API based on eprinttools Go package",
    long_description = read(readme_txt),
    author = author,
    author_email = author_email,
    url = "https://caltechlibrary.github.io/eprinttools",
    download_url = "https://github.com/caltechlibrary/eprinttools/latest/releases",
    license = meta["license"],
    packages = ["eprinttools"],
    data_files = [
        (site_package_location, [os.path.join("eprinttools", shared_library_name)]),
    ],
    platforms = [platform],
    keywords = ["API", "library", "repository"],
    classifiers = [
        "Development Status :: Alpha",
        "Environment :: Console",
        "Programming Language :: Python",
        "Programming Language :: Python :: 3",
        "Programming Language :: Other",
        "Programming Language :: Go",
        "Intended Audience :: Science/Research",
        "License :: OSI Approved :: BSD License",
        OS_Classifier
    ]
)
