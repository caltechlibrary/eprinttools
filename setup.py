#!/usr/bin/env python3

#from distutils.core import setup
#from site import getsitepackages
#site_package_location = os.path.join(getsitepackages()[0], "dataset")

from setuptools import setup, find_packages, Command
from shutil import rmtree, copyfile

import sys
import os,io
import json

readme_md = "README.md"
readme_txt = "README.txt"
codemeta_json = "codemeta.json"

def read(fname):
    with open(fname, mode = "r", encoding = "utf-8") as f:
        src = f.read()
    return src

if os.path.exists(codemeta_json) == False:
    copyfile(os.path.join("..", codemeta_json), codemeta_json)

# Let's pickup as much metadata as we can from codemeta.json
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
        author = given + " " + family
    else:
        author = author + ", " + given + " " + family
    if len(author_email) == 0:
        author_email = email
    else:
        author_email = author_email + ", " + email
description = meta['description']
keywords = meta['keywords']
url = meta['codeRepository']
download = meta['downloadUrl']
license = meta['license']
name = meta['name']

REQUIRES_PYTHON = '>=3.7.0'

# What packages are required for this module to be executed?
REQUIRED = [ "py_dataset" ]

# What packages are optional?
EXTRAS = {}

# The rest you shouldn't have to touch too much :)
# ------------------------------------------------
# Except, perhaps the License and Trove Classifiers!
# If you do change the License, remember to change the Trove Classifier for that!

here = os.path.abspath(os.path.dirname(__file__))

# Import the README and use it as the long-description.
# Note: this will only work if 'README.md' is present in your MANIFEST.in file!
try:
    with io.open(os.path.join(here, readme_md), encoding='utf-8') as f:
        long_description = '\n' + f.read()
except FileNotFoundError:
    long_description = description

class UploadCommand(Command):
    """Support setup.py upload."""

    description = 'Build and publish the package.'
    user_options = []

    @staticmethod
    def status(s):
        """Prints things in bold."""
        print('\033[1m{0}\033[0m'.format(s))

    def initialize_options(self):
        pass

    def finalize_options(self):
        pass

    def run(self):
        try:
            self.status('Removing previous builds…')
            rmtree(os.path.join(here, 'dist'))
        except OSError:
            pass

        self.status('Building Source and Wheel distribution…')
        os.system('{0} setup.py sdist bdist_wheel '.format(sys.executable))

        self.status('Uploading the package to PyPI via Twine…')
        os.system('twine upload dist/*')

        sys.exit()

# Now that we know everything configure out setup
setup(name = name,
    version = version,
    description = description,
    long_description=long_description,
    long_description_content_type='text/markdown',
    python_requires=REQUIRES_PYTHON,
    install_requires=REQUIRED,
    extras_require=EXTRAS,
    author = author,
    author_email = author_email,
    url = url,
    download_url = download,
    license = license,
    packages = find_packages(exclude=["*.tests", "*.tests.*", "tests.*", "tests", "*_test.py"]),
    package_data={},
    keywords = keywords,
    classifiers = [
        "Development Status :: 4 - Beta",
        "Environment :: Console",
        "Programming Language :: Python",
        "Programming Language :: Python :: 3",
        "Programming Language :: Other",
        "Intended Audience :: Science/Research",
        "Topic :: Scientific/Engineering",
        "License :: OSI Approved :: BSD License",
        "Operating System :: Microsoft :: Windows :: Windows 10",
        "Operating System :: POSIX :: Linux",
        "Operating System :: MacOS :: MacOS X"
    ],
    cmdclass={
        'upload': UploadCommand,
    }
)
