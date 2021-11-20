'''
  cltests is a lightweight Test library for checking type and values
  for equality by implementing simple boolean test functions.
'''
import sys

from .T import T
from .TestSet import TestSet


def IsSuccessful(result):
    if result == True:
        sys.exit(0)
    else:
        sys.exit(1)
