Z3_DIR  ?= ../z3
Z3_INC  := $(Z3_DIR)/src/api
Z3_LIB  := $(Z3_DIR)/build/libz3.so
Z3_RPATH := $(Z3_DIR)/build

CGO_ENV := CGO_CFLAGS="-I$(Z3_INC)" \
					 CGO_LDFLAGS="$(Z3_LIB) -lpthread" \
					 CGO_LDFLAGS_ALLOW=".*"

.PHONY: build build-z3 test clean

# Build every package, including internal/encoder, internal/solver, and cmd/.
build:
	$(CGO_ENV) go build ./...

# Build just the Z3 smoke test binary.
build-z3:
	$(CGO_ENV) go build -ldflags="-r $(Z3_RPATH)" -o bin/z3 cmd/z3/main.go

bin-verify:
	$(CGO_ENV) go build -ldflags="-r $(Z3_RPATH)" -o bin/verify ./cmd/verify

run-verify: bin-verify
	./bin/verify --data testdata/interfaces-pass.json --mtu-min 1500
	./bin/verify --data testdata/interfaces-fail.json --mtu-min 1500 || true

test:
	$(CGO_ENV) go test ./...

clean:
	rm -rf bin
