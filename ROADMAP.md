
# Roadmap

1. [x] Dropped dataset integeration in Golang codebase
2. [x] Release v0.1.0 as set of standalone tools that are easily wrapped by Python
3. [x] Implement lightweight wrapper object, eprints3x in python that uses `eputil` to interact with EPrints and format the JSON version of eprints metadata
4. [x] Implement harvester.py using eprints3x library that replicates the functionality of the old `ep` go based cli
    + renamed to demo-harvester-full.py and demo-harvester-recent.py
5. Convert convert eprints3x codebase to pure Python 3.x replacing calls to `eputil`
6. Convert remaining Go based cli to Python3 and compile to standard alone cli where appropriate

