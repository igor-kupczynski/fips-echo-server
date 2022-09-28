# fips-echo-server

Experiments with FIPS-compliant Golang crypto

This repo has a companion blog post, check https://kupczynski.info/posts/fips-golang/

 ## Http echo server

We have a simple http(s) echo server. It echos up to 140 characters of whatever you send it to it.

You can test it locally (assuming you have go toolchain installed).

Let's build and run it (press `Ctrl+C` to stop):
```sh
$ go build .
$ ./fips-echo-server
2022/09/28 22:53:25 Listening on https://localhost:8443 with cert=certs/domain.pem and key=certs/domain.key
```

And then test it in another terminal:
```sh
$ curl --cacert certs/ca.pem https://localhost:8443 -d "hello"
hello
```

Note that I'm embedding some self-signed certs in the `certs` folder. We need `--cacert certs/ca.pem` flag for curl to
trust them.

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

It is based on `golang`, which is the _official_ go toolchain image. Iit adds the project folder to `/app` in
the container and compiles it into a binary. Finally, it exposes the `:8443` port and sets the produced binary as
the default startup command for the container.

We have a [`Makefile`](./Makefile) to save us some typing. Check it out for the exact commands.

We can build the docker container:
```sh
$ make build
[+] Building 1.0s (10/10) FINISHED
..
 => => naming to docker.io/library/fips-echo-server:go1.18.6
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
2022/09/28 20:54:35 Listening on https://0.0.0.0:8443 with cert=certs/domain.pem and key=certs/domain.key
```

The main motivation for the dockerized build & run in the repo is to demonstrate the FIPS compliant go toolchain.
It will save us the hassle of setting it latter.

## TLS setup

Suggested reading from Cloudflare on the subject —
[_So you want to expose Go on the Internet_](https://blog.cloudflare.com/exposing-go-on-the-internet/).

What are the ciphers we use by default? To test that we'll use the [`testssl.sh`](https://github.com/drwetter/testssl.sh)
script.

```sh
# not in the project directory
$ git clone git@github.com:drwetter/testssl.sh.git
$ cd testssl.sh
$ ./testssl.sh localhost:8443
...
 Testing server's cipher preferences 

Hexcode  Cipher Suite Name (OpenSSL)       KeyExch.   Encryption  Bits     Cipher Suite Name (IANA/RFC)
-----------------------------------------------------------------------------------------------------------------------------
SSLv2
 - 
SSLv3
 - 
TLSv1 (server order)
 xc013   ECDHE-RSA-AES128-SHA              ECDH 521   AES         128      TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA                 
 xc014   ECDHE-RSA-AES256-SHA              ECDH 521   AES         256      TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA                 
 x2f     AES128-SHA                        RSA        AES         128      TLS_RSA_WITH_AES_128_CBC_SHA                       
 x35     AES256-SHA                        RSA        AES         256      TLS_RSA_WITH_AES_256_CBC_SHA                       
 xc012   ECDHE-RSA-DES-CBC3-SHA            ECDH 521   3DES        168      TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA                
 x0a     DES-CBC3-SHA                      RSA        3DES        168      TLS_RSA_WITH_3DES_EDE_CBC_SHA                      
TLSv1.1 (server order)
 xc013   ECDHE-RSA-AES128-SHA              ECDH 521   AES         128      TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA                 
 xc014   ECDHE-RSA-AES256-SHA              ECDH 521   AES         256      TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA                 
 x2f     AES128-SHA                        RSA        AES         128      TLS_RSA_WITH_AES_128_CBC_SHA                       
 x35     AES256-SHA                        RSA        AES         256      TLS_RSA_WITH_AES_256_CBC_SHA                       
 xc012   ECDHE-RSA-DES-CBC3-SHA            ECDH 521   3DES        168      TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA                
 x0a     DES-CBC3-SHA                      RSA        3DES        168      TLS_RSA_WITH_3DES_EDE_CBC_SHA                      
TLSv1.2 (server order -- server prioritizes ChaCha ciphers when preferred by clients)
 xc02f   ECDHE-RSA-AES128-GCM-SHA256       ECDH 521   AESGCM      128      TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256              
 xc030   ECDHE-RSA-AES256-GCM-SHA384       ECDH 521   AESGCM      256      TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384              
 xcca8   ECDHE-RSA-CHACHA20-POLY1305       ECDH 521   ChaCha20    256      TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256        
 xc013   ECDHE-RSA-AES128-SHA              ECDH 521   AES         128      TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA                 
 xc014   ECDHE-RSA-AES256-SHA              ECDH 521   AES         256      TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA                 
 x9c     AES128-GCM-SHA256                 RSA        AESGCM      128      TLS_RSA_WITH_AES_128_GCM_SHA256                    
 x9d     AES256-GCM-SHA384                 RSA        AESGCM      256      TLS_RSA_WITH_AES_256_GCM_SHA384                    
 x2f     AES128-SHA                        RSA        AES         128      TLS_RSA_WITH_AES_128_CBC_SHA                       
 x35     AES256-SHA                        RSA        AES         256      TLS_RSA_WITH_AES_256_CBC_SHA                       
 xc012   ECDHE-RSA-DES-CBC3-SHA            ECDH 521   3DES        168      TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA                
 x0a     DES-CBC3-SHA                      RSA        3DES        168      TLS_RSA_WITH_3DES_EDE_CBC_SHA                      
TLSv1.3 (no server order, thus listed by strength)
 x1302   TLS_AES_256_GCM_SHA384            ECDH 253   AESGCM      256      TLS_AES_256_GCM_SHA384                             
 x1303   TLS_CHACHA20_POLY1305_SHA256      ECDH 253   ChaCha20    256      TLS_CHACHA20_POLY1305_SHA256                       
 x1301   TLS_AES_128_GCM_SHA256            ECDH 253   AESGCM      128      TLS_AES_128_GCM_SHA256                             

 Has server cipher order?     yes (OK) -- only for < TLS 1.3
 Negotiated protocol          TLSv1.3
 Negotiated cipher            TLS_AES_128_GCM_SHA256, 253 bit ECDH (X25519)
...
```

`testssl.sh` presents a long report, but for us the important part is given above. By default, go 1.18 supports
TLSv1.0--TLSv1.3.

Let's be more strict here and select only the
[protocols and ciphers recommended by Mozilla for a _modern_ configuration](https://wiki.mozilla.org/Security/Server_Side_TLS).

_You can also add the args to `CMD` in `Dockerfile`_.

```sh
$ ./fips-echo-server -tlsVersion TLSv1.3 -tlsCiphers TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:ECDHE-RSA-AES128-GCM-SHA256
```

This results is:
```sh
...
 Testing server's cipher preferences 

Hexcode  Cipher Suite Name (OpenSSL)       KeyExch.   Encryption  Bits     Cipher Suite Name (IANA/RFC)
-----------------------------------------------------------------------------------------------------------------------------
SSLv2
 - 
SSLv3
 - 
TLSv1
 - 
TLSv1.1
 - 
TLSv1.2
 - 
TLSv1.3 (no server order, thus listed by strength)
 x1302   TLS_AES_256_GCM_SHA384            ECDH 253   AESGCM      256      TLS_AES_256_GCM_SHA384                             
 x1303   TLS_CHACHA20_POLY1305_SHA256      ECDH 253   ChaCha20    256      TLS_CHACHA20_POLY1305_SHA256                       
 x1301   TLS_AES_128_GCM_SHA256            ECDH 253   AESGCM      128      TLS_AES_128_GCM_SHA256                             

 Has server cipher order?     no (TLS 1.3 only)
 Negotiated protocol          TLSv1.3
 Negotiated cipher            TLS_AES_128_GCM_SHA256, 253 bit ECDH (X25519) (limited sense as client will pick)
```

## FIPS compliant version

Please check the
[`boringcrypto-1.18` branch](https://github.com/igor-kupczynski/fips-echo-server/compare/main-1.18...boringcrypto-1.18)
branch for details. Compare it with the current one to find the changes needed to support a FIPS mode in our app.