# fips-echo-server

Experiments with FIPS-compliant Golang crypto

 ## Http echo server

We have a simple http(s) echo server. It echos up to 140 characters of whatever you send it to it.

Let's run it (press `Ctrl+C` to stop):
```sh
$ make build
$ ./fips-echo-server
2019/11/29 19:26:00 Listening on https://localhost:8443 with cert=certs/domain.pem and key=certs/domain.key
```

And then in another terminal:
```sh
$ curl --cacert certs/ca.pem https://localhost:8443 -d "hello"
hello
```

Note that I'm embedding some self-signed certs in the `certs` folder. We need `--cacert certs/ca.pem` flag for curl to trust them.

You can also run the tests:
```sh
$ make test
?   	github.com/igor-kupczynski/fips-echo-server	[no test files]
=== RUN   TestServe
=== RUN   TestServe/Echo_the_message_back_to_the_client
=== RUN   TestServe/Limit_to_140_characters
--- PASS: TestServe (0.04s)
    --- PASS: TestServe/Echo_the_message_back_to_the_client (0.02s)
    --- PASS: TestServe/Limit_to_140_characters (0.02s)
PASS
ok  	github.com/igor-kupczynski/fips-echo-server/echo	0.222s
```