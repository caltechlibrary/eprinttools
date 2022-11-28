Repositories (end point)
========================

This end point lists the repositories known to the __ep3apid__ service.

- '/repositories' returns a JSON array of repository names defined in settings.json
- '/repositories/help' returns this documentation.

Example
-------

In this example we assume the __ep3apid__ services is running on "localhost:8484" and is configured to support two repositories "lemurprints" and "test3". We are using curl to retrieve the data.

```shell
    curl -X GET http://localhost:8484/repositories
```

This should return a JSON array like

```json
    [
        "lemurprints",
        "test3"
    ]
```




