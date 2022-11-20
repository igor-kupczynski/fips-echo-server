# Start with an official image
FROM golang:1.19.3

# Bundle the project source in the container
RUN mkdir -p /app
ADD . /app
WORKDIR /app

# Build a binary and assert that it uses boringcrypto instead of the native golang crypto
RUN GOEXPERIMENT=boringcrypto go build . && \
    go tool nm fips-echo-server > tags.txt && \
    grep '_Cfunc__goboringcrypto_' tags.txt 1> /dev/null

# Expose a port and set the default run command for the container
EXPOSE 8443
CMD [ "./fips-echo-server", "-address", "0.0.0.0:8443", "-fipsMode", "true" ]