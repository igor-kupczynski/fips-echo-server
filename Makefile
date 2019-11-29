label = fips-echo-server:latest
addr = 127.0.0.1:8443

help:
	@ echo "Goals:"
	@ echo "  build - build the echo server docker container and tag it as $(label)"
	@ echo "  test  - run the tests in the container"
	@ echo "  run   - run the server in the container and expose it to $(addr)"

build:
	@ docker build -t $(label) .

test:
	@ docker run --rm $(label) go test -v ./...

run:
	@ docker run --rm -p $(addr):8443 $(label)