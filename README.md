# fips-echo-server

Experiments with FIPS-compliant Golang crypto

 ## Http echo server

We have a simple http(s) echo server. It echos up to 140 characters of whatever you send it to it.

You can test it locally (assuming you have go toolchain installed).

Let's build and run it (press `Ctrl+C` to stop):
```sh
$ go build .
$ ./fips-echo-server
2019/11/29 19:26:00 Listening on https://localhost:8443 with cert=certs/domain.pem and key=certs/domain.key
```

And then test it in another terminal:
```sh
$ curl --cacert certs/ca.pem https://localhost:8443 -d "hello"
hello
```

Note that I'm embedding some self-signed certs in the `certs` folder. We need `--cacert certs/ca.pem` flag for curl to trust them.

You can also run the tests:
```sh
$ go test ./...
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

## Dockerized version

To abstract away the go toolchain you may have installed locally we also provide a [Dockerfile](./Dockerfile). It starts with `golang`, which is the _official_ go toolchain image. The it adds the project folder to `/app` in the container and compiles it into a binary. Finally, it exposes the `:8443` port and sets the produced binary as the default startup command for the container.

We have a [`Makefile`](./Makefile) to save us some typing. Check it out for the exact commands.

We can build the docker container:
```sh
$ make build
Sending build context to Docker daemon  7.683MB
(...)
Successfully tagged fips-echo-server:latest
```

Run the tests within the container:
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
ok  	github.com/igor-kupczynski/fips-echo-server/echo	0.045s
```

And finally run the container with the echo server:
```sh
$ make run
2019/11/29 23:40:46 Listening on https://0.0.0.0:8443 with cert=certs/domain.pem and key=certs/domain.key
```

Since golang is multiplatform, docker may seem like an overkill. There are  some advantages, e.g. when it comes to CI pipeline or artifact distribution. The main motivation for the dockerized build & run in the repo is to demonstrate the FIPS compliant go toolchain. It will save us the hassle of setting it latter.