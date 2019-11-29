help:
	@ echo "Goals:"
	@ echo "  build - build the server"
	@ echo "  test  - run the tests"

build:
	@ go build .

test:
	@ go test -v ./...