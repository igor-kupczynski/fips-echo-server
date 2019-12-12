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

To build/run independently of a local go toolchain we provide a [Dockerfile](./Dockerfile).

It is based on `golang`, which is the _official_ go toolchain image. The it adds the project folder to `/app` in the container and compiles it into a binary. Finally, it exposes the `:8443` port and sets the produced binary as the default startup command for the container.

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

The main motivation for the dockerized build & run in the repo is to demonstrate the FIPS compliant go toolchain. It will save us the hassle of setting it latter.

## TLS setup

Suggested reading from Cloudflare on the subject — [_So you want to expose Go on the Internet_](https://blog.cloudflare.com/exposing-go-on-the-internet/).

What are the ciphers we use by default? To test that we'll use the [`testssl.sh`](https://github.com/drwetter/testssl.sh) script.
```sh
# not in the project directory
$ git clone git@github.com:drwetter/testssl.sh.git
$ cd testssl.sh
$ ./testssl.sh localhost:8443
...
 Testing server preferences

 Has server cipher order?     yes (OK) -- TLS 1.3 and below
 Negotiated protocol          TLSv1.3
 Negotiated cipher            TLS_AES_128_GCM_SHA256, 253 bit ECDH (X25519)
 Cipher order
    TLSv1:     ECDHE-RSA-AES128-SHA ECDHE-RSA-AES256-SHA AES128-SHA AES256-SHA ECDHE-RSA-DES-CBC3-SHA DES-CBC3-SHA
    TLSv1.1:   ECDHE-RSA-AES128-SHA ECDHE-RSA-AES256-SHA AES128-SHA AES256-SHA ECDHE-RSA-DES-CBC3-SHA DES-CBC3-SHA
    TLSv1.2:   ECDHE-RSA-AES128-GCM-SHA256 ECDHE-RSA-AES256-GCM-SHA384 ECDHE-RSA-CHACHA20-POLY1305 ECDHE-RSA-AES128-SHA ECDHE-RSA-AES256-SHA AES128-GCM-SHA256 AES256-GCM-SHA384
               AES128-SHA AES256-SHA ECDHE-RSA-DES-CBC3-SHA DES-CBC3-SHA
    TLSv1.3:   TLS_AES_128_GCM_SHA256 TLS_CHACHA20_POLY1305_SHA256 TLS_AES_256_GCM_SHA384
...
```

`testssl.sh` presents a long report, but for us the important part is given above. By default go 1.13 support TLSv1.0--TLSv1.3. Let's be more strict here and select only the [protocols and ciphers recommended by Mozilla for a _modern_ configuration](https://wiki.mozilla.org/Security/Server_Side_TLS).

_You can also add the args to `CMD` in `Dockerfile`_.

```sh
$ ./fips-echo-server -tlsVersion TLSv1.3 -tlsCiphers TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:ECDHE-RSA-AES128-GCM-SHA256
```

This results is:
```sh
...
 Testing server preferences

 Has server cipher order?     yes (TLS 1.3 only)
 Negotiated protocol          TLSv1.3
 Negotiated cipher            TLS_AES_128_GCM_SHA256, 253 bit ECDH (X25519)
 Cipher order
    TLSv1.3:   TLS_AES_128_GCM_SHA256 TLS_CHACHA20_POLY1305_SHA256 TLS_AES_256_GCM_SHA384
```

## FIPS compliant version

Please check the [`boringcrypto` branch](tree/boringcrypto) branch for details. Compare it with the current one to find the changes needed to support a FIPS mode in our app.