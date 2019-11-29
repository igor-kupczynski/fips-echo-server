# fips-echo-server

Experiments with FIPS-compliant Golang crypto

 ## Http echo server

We have a simple http(s) echo server. It echos up to 140 characters of whatever you send it to it.

Let's run it:
```sh
$ go run . -port 8443
2019/11/29 18:49:04 Listening on :8443
```

And then in another terminal:
```sh
$ curl http://localhost:8443 -s -d "hello"
hello
```

You can also run the tests:
```sh
$ cd echo     
$ go test
PASS
ok      github.com/igor-kupczynski/fips-echo-server/echo        0.369s
```