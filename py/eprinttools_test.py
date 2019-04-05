#!/usr/bin/env python3
import os
import sys
import eprinttools
import random
import datetime

#
# Tests
#
def test_get_keys(t, eprint_url, auth_type = "", username = "", secret = ""):
    #t.verbose_on() # turn verboseness on for debugging
    cfg = eprinttools.cfg(eprint_url, auth_type, username, secret)
    keys = eprinttools.get_keys(cfg)
    if len(keys) == 0:
        t.error(f"Expected more than zero keys for {cfg['eprint_url']}")
        t.fail_now()
    else:
        t.print("found", len(keys), f"for {cfg}")


def test_get_modified_keys(t, eprint_url, auth_type = "", username = "", secret = ""):
    #t.verbose_on() # turn verboseness on for debugging
    cfg = eprinttools.cfg(eprint_url, auth_type, username, secret)
    # we are checking to see if we have recently modified keys (last 30 days)
    end = datetime.datetime.now()
    start = end - datetime.timedelta(days = 30)
    t.verbose_on()
    eprinttools.verbose_on()
    t.print(f"Checking for {eprint_url} (username: {username}, auth: {auth_type}) datetime range {start} to {end} (this can take a while)")
    keys = eprinttools.get_modified_keys(cfg, start, end)
    if keys == None or len(keys) == 0:
        t.error(f"expected more than zero keys for get_modified_keys({cfg['eprint_url']}, {start}, {end}")
        t.fail_now()
    else:
        t.print("found", len(keys), f"for {cfg['url']}")
    eprinttools.verbose_off()


def test_get_metadata(t, eprint_url, auth_type = 0, username = "", secret = ""):
    test_name = t.test_name()
    t.verbose_on() # turn verboseness on for debugging
    collection_name = test_name + ".ds"
    cfg = eprinttools.cfg(eprint_url, auth_type, username, secret, collection_name)
    keys = eprinttools.get_keys(cfg)
    if len(keys) == 0:
        t.error(f"Can't test {test_name} without keys, got zero keys")
        t.fail_now()
        return
    #NOTE: Pick some random keys to test getting metadata records.
    # So we don't go through really large collections.
    collection_keys = []
    check_keys = []
    for i in range(100):
        key = random.choice(keys)
        if key not in check_keys:
            check_keys.append(key)
        if len(check_keys) > 50:
            break
    for key in check_keys:
        # We are going to try to get the metadata for the EPrint record but not store it in a dataset collectin...
        data = eprinttools.get_metadata(cfg, key, False)
        e_msg = eprinttools.error_message()
        if len(data) == 0 or e_msg != "":
            if e_msg.startswith("401"):
                t.print(f"found {key}, requires authentication")
            elif "buffer" in e_msg or "deletion" in e_msg or "indox" in e_msg:
                t.print(f"skipping {key}, non-public status")
            else:
                t.error(f"Expected data for {key}, got {data} '{e_msg}'")
                t.fail_now()
        else:
            t.print(f"found {key} with data")
            collection_keys.append(key)

    # Check to see if we can retrieved the buffered XML
    for key in collection_keys:
        data = eprinttools.get_metadata(cfg, key, True)
        e_msg = eprinttools.error_message()
        xml_src = eprinttools.get_buffered_xml()
        if len(data) == 0 or e_msg != "":
            e_msg = eprinttools.error_message()
            if e_msg.startswith("401") == False:
                t.error(f"Expected data for {key}, got {data} {e_msg}")
                t.fail_now()
            else:
                t.print(f"found {key}, requires authentication")
        else:
            t.print(f"found {key} with data")
            if len(xml_src) == 0:
                t.error("Could not get xml buffer contents for {key}")


#
# Test harness
#
class ATest:
    def __init__(self, test_name, verbose = False):
        self.pid = os.getpid()
        self._test_name = test_name
        self._error_count = 0
        self._verbose = verbose

    def test_name(self):
        return self._test_name

    def is_verbose(self):
        return self._verbose

    def verbose_on(self):
        self._verbose = True
       
    def verbose_off(self):
        self._verbose = False

    def print(self, *msg):
        if self._verbose == True:
            print(f"(pid: {self.pid})", *msg)

    def error(self, *msg):
        fn_name = self._test_name
        self._error_count += 1
        print(f"(pid: {self.pid}) {fn_name} failed, ", *msg)

    def error_count(self):
        return self._error_count

    def fail_now(self):
        fn_name = self._test_name
        error_count = self._error_count
        print(f"(pid: {self.pid}) {fn_name} failed, {error_count}")
        sys.exit(1)


class TestRunner:
    def __init__(self, set_name):
        self._set_name = set_name
        self._tests = []
        self._error_count = 0


    def add(self, fn, params = []):
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
        print("PASS")
        print("Ok", __file__)
        sys.exit(0)


def setup():
    eprint_url, auth_type, username, secret = "", "", "", ""
    conf = eprinttools.envcfg()
    if 'url' in conf:
        eprint_url = conf['url']
    if 'auth_type' in conf:
        auth_type = conf['auth_type']
    if 'username' in conf:
        username = conf['username']
    if 'secret' in conf:
        secret = conf['secret']
    return eprint_url, auth_type, username, secret

#
# Run tests
#
if __name__ == "__main__":
    version = eprinttools.version()
    eprint_url, auth_type, username, secret = setup()
    if eprint_url == "":
        print(f"Skipping tests for eprinttools {version}, EPRINT_URL not set in the environment")
        sys.exit(0)
    if username != "" and auth_type == "":
        print(f"Skipping tests for eprinttools {version}, EPRINT_AUTH_TYPE not set in the environment")
        sys.exit(1)
    test_runner = TestRunner(os.path.basename(__file__))
    test_runner.add(test_get_keys, [eprint_url, auth_type, username, secret])
    if "-quick" in sys.argv or "quick" in sys.argv:
        print("Skipping test_get_modified_keys, -quick option used")
    else:
        test_runner.add(test_get_modified_keys, [eprint_url, auth_type, username, secret])
    test_runner.add(test_get_metadata, [eprint_url, auth_type, username, secret])
    test_runner.run()

