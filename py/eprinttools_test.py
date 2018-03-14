#!/usr/bin/env python3
import os
import sys
import eprinttools

#
# Tests
#
def test_get_keys(t, eprint_url, auth_type = 0, username = "", secret = ""):
    t.errorf("not implemented")

def test_get_metadata(t, eprint_url, auth_type = 0, username = "", secret = ""):
    t.errorf("not implemented")


#
# Test harness
#
class ATest:
    def __init__(self, test_name):
        self._test_name = test_name
        self._error_count = 0

    def errorf(self, msg):
        fn_name = self._test_name
        self._error_count += 1
        print(f"\t{fn_name}", msg)

    def error_count(self):
        return self._error_count

class TestRunner:
    def __init__(self, set_name):
        self._set_name = set_name
        self._tests = []
        self._error_count = 0

    def add(self, fn, params):
        self._tests.append((fn, params))

    def run(self):
        for test in self._tests:
            fn_name = test[0].__name__
            t = ATest(fn_name)
            fn, params = test[0], test[1]
            fn(t, *params)
            error_count = t.error_count()
            if error_count > 0:
                print(f"\t\t{fn_name} failed, {error_count} errors found")
            self._error_count += error_count
        error_count = self._error_count
        set_name = self._set_name
        if error_count > 0:
            print(f"Failed {set_name}, {error_count} total errors found")
            sys.exit(1)
        print("Success!")
        sys.exit(0)


def setup():
    eprint_url = os.getenv("EPRINT_URL")
    auth_type = os.getenv("EPRINT_AUTH_TYPE")
    username = os.getenv("EPRINT_USER")
    secret = os.getenv("EPRINT_PASSWD")
    if auth_type == None:
        auth_type = ""
    if username == None:
        username = ""
    if secret == None:
        secret = ""
    return eprint_url, auth_type, username, secret

#
# Run tests
#
if __name__ == "__main__":
    version = eprinttools.version()
    eprint_url, auth_type, username, secret = setup()
    if eprint_url == None:
        print(f"Skipping tests for eprinttools {version}, EPRINT_URL not set in the environment")
        sys.exit(0)
    print(f"Testing eprinttools {version}")
    test_runner = TestRunner(os.path.basename(__file__))
    test_runner.add(test_get_keys, [eprint_url, auth_type, username, secret])
    test_runner.add(test_get_metadata, [eprint_url, auth_type, username, secret])
    test_runner.run()

