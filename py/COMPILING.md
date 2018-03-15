
# How to compile libeprint3 shared library

## Go Setup

I am using Go v1.10 and have cgo available. 

## Python

This package is compiled for Python 3.6 and the python code 
uses 3.6 features (e.g. f strings)

## Mac OS X

(didn't seem to make a difference, rsd 2018-03-15)

You need to have 'pkg-config: python-3.6' installed. If you are
using Mac Ports

```shell
    sudo port install py36-pkgconfig
    export PKG_CONFIG_PATH="/opt/local/Library/Frameworks/Python.framework/Versions/3.6/lib/pkgconfig"
```
