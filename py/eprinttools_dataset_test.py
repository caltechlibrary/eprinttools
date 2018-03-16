#!/usr/bin/env python3
import os
import shutil
import sys
import dataset
import eprinttools
import random
import datetime

#
# Tests
#
def test_get_metadata(t, eprint_url, auth_type, username, secret, collection_name):
    if os.path.exists(collection_name):
        shutil.rmtree(collection_name)
    ok = dataset.init(collection_name)
    if ok == False:
        t.error(f"Can't initialize {collection_name}")
        return
    t.verbose_off() # turn verboseness on for debugging
    test_name = t.test_name()
    cfg = eprinttools.cfg(eprint_url, auth_type, username, secret, collection_name)
    keys = eprinttools.get_keys(cfg)
    if len(keys) == 0:
        t.error(f"Can't test {test_name} without keys, got zero keys")
        return
    #FIXME: Pick some random keys to test getting metadata records!
    collection_keys = []
    check_keys = []
    for i in range(100):
        key = random.choice(keys)
        if key not in check_keys:
            check_keys.append(key)
        if len(check_keys) > 50:
            break
    t.print(f"Calculating the keys in sample that will get stored in the collection {collection_name}")
    for key in check_keys:
        # We are going to try to get the metadata for the EPrint record but not store it in a dataset collectin...
        data = eprinttools.get_metadata(cfg, key, False)
        if len(data) == 0:
            e_msg = eprinttools.error_message()
            if e_msg.startswith("401") == False:
                t.error(f"Expected data for {key}, got {data}")
            else:
                t.print(f"found {key}, requires authentication")
        else:
            t.print(f"found {key} with data")
            collection_keys.append(key)
            dataset.create(collection_name, key, data)

    eprinttools.verbose_on() # DEBUG
    dataset.verbose_on() # DEBUG
    key_cnt = len(collection_keys)
    t.print(f"harvesting {key_cnt} keys to {collection_name}")
    dataset.verbose_on() # DEBUG
    t.verbose_on() # DEBUG
    for key in collection_keys:
        data = eprinttools.get_metadata(cfg, key, True)
        xml_src = eprinttools.get_buffered_xml()
        e_msg = eprinttools.error_message()
        if e_msg != "" and e_msg.startswith("401") == False:
                t.error(f"Expected data for {key}, got {e_msg}")
        elif e_msg != "":
                t.print(f"Skipped {key}, requires authentication")
        elif len(data) == 0:
            t.print(f"warning: {key} did not return any data")
        else:
            t.print(f"found {key} with data {cfg}")
            rec = dataset.read(collection_name, key)
            e_msg = dataset.error_message()
            if len(rec) == 0 or e_msg != "":
                t.error(f"Should be able to read '{key}' in '{collection_name}' {e_msg}")
                return
            attachment_name = f"{key}.xml"
            with open(attachment_name, mode = "w", encoding = "utf-8") as f:
                f.write(xml_src)
            ok = dataset.attach(collection_name, key, [attachment_name])
            if ok == False:
                t.error("Count not attach {attachment_name} to {key} in {collection_name}")
                return
            os.remove(attachment_name)        

    keys = dataset.keys(collection_name)
    if len(keys) != len(collection_keys):
        t.error("expected collection keys to match batch harvested")

#
# Test harness
#
class ATest:
    def __init__(self, test_name, verbose = False):
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
            print(*msg)

    def error(self, *msg):
        fn_name = self._test_name
        self._error_count += 1
        print(f"{fn_name} failed, ", *msg)

    def error_count(self):
        return self._error_count

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
    ep_version = eprinttools.version()
    ds_version = dataset.version()

    eprint_url = os.getenv("EPRINT_URL")
    auth_type = os.getenv("EPRINT_AUTH_TYPE")
    username = os.getenv("EPRINT_USER")
    secret = os.getenv("EPRINT_PASSWD")
    collection_name = "test_get_metadata.ds" #os.getenv("DATASET")

    if eprint_url == None:
        print(f"Skipping tests for eprinttools {ep_version}, EPRINT_URL not set in the environment")
        sys.exit(1)
    if collection_name == None:
        print(f"Skipping tests for eprinttools with dataset {ds_version}, DATASET not set in the environment")
        sys.exit(1)
    if os.path.exists(collection_name) == False:
        ok = dataset.init(collection_name)
        if ok == False:
            print(f"Could not init {collection_name}")
            sys.exit(1)


    if auth_type == None:
        auth_type = ""
    if username == None:
        username = ""
    if secret == None:
        secret = ""
    return eprint_url, auth_type, username, secret, collection_name

#
# Run tests
#
if __name__ == "__main__":
    eprint_url, auth_type, username, secret, collection_name = setup()
    test_runner = TestRunner(os.path.basename(__file__))
    test_runner.add(test_get_metadata, [eprint_url, auth_type, username, secret, collection_name])
    test_runner.run()

