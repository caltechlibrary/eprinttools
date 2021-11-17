Test Setup
==========

It is important to be able to test the code that makes up eprinttools. EPrints is complex and nuanced and testing helps insure that our code is addressing that. For running the test suite you need to have MySQL setup with a test database and one or more EPrints repository available and setup for REST access. The REST client testing is read only and can be skipped if you're not using it. You need at least one MySQL database setup using the same structure that is present in our EPrints 3.3.16 repositories. You also need a "test-settings.json" file for the configured test.

The test database used to test both SQL level interaction (which includes writes and deletes) should be named "lemurprints". It should now have records in it when you start. The Scheme for that database is found in the directory `srctest/lemurprints-setup-schema.sql`. Assuming your MySQL client is configured for loading databased you can run

~~~{.bash}
    mysql --execute 'CREATE DATABASE IF NOT EXIST lemurprints'
    mysql lemurprints < srctest/lemurprints-setup-schema.sql
~~~

You're `test-settings.json` file should look something like this (replace the text in capital letters appropriately).

~~~{.json}
{
    "logfile": "eprinttools-test.log",
    "repositories": {
        "lemurprints": {
            "dsn": "USERNAME:PASSWORD@/lemurprints",
            "base_url": "http://lemurprints.example.edu",
            "rest": "https://USERNAME:PASSWORD@TEST_HOSTNAME_FOR_REST_CLIENT",
            "write": true
        }
    }
}
~~~

You can omit the "rest" key/value if you're not going to test the REST client.

When running tests any log output will be sent to the "logfile" value.  If that is not set then you'll have at least one test failure (the one testing if log is being sent to a file).

The test database "lemurprints" need to be "write" enabled. If you add additional repository databases DO NOT enable "write" as the test sequence will attempt to clear the database before each test run.



