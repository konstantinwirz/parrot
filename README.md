# parrot
[![Build Status](https://travis-ci.com/konstantinwirz/parrot.svg?branch=master)](https://travis-ci.com/konstantinwirz/parrot)

A very simple mock http server, that echoes created resources back like a parrot

## how to start the server
```
> docker run -d -p 8080:8080 kitov/parrot parrot
```

## examples

### create a resource
```
> http POST :8080/foo/25 a=b c:=2 d:='{"f":"g"}'

HTTP/1.1 201 Created
Content-Length: 32
Content-Type: application/json
Date: Sun, 28 Apr 2019 17:54:28 GMT
Server: Parrot (go1.12.4)

{
    "code": 201,
    "message": "created"
}
```

### request a resource
```
> http GET :8080/foo/25

HTTP/1.1 200 OK
Content-Length: 35
Content-Type: application/json
Date: Sun, 28 Apr 2019 17:53:51 GMT
Server: Parrot (go1.12.4)

{
    "a": "b",
    "c": 2,
    "d": {
        "f": "g"
    }
}

```