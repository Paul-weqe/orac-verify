# orac-verify

Encodes normalized network data as SMT constraints and checks them with Z3.
Takes JSON produced by [orac-data](https://github.com/Paul-weqe/orac-data),
builds an assertion tree, solves it, and reports whether a rule holds.

Current capability: interface MTU verification. Given a JSON file of
normalized interface data and a minimum MTU, the tool reports whether every
interface meets the threshold. On violation, it prints the offending values
and exits nonzero.

## Requirements

- Go 1.25 or newer
- A C toolchain (gcc or clang)
- Z3 built from source with Go bindings and shared library enabled
- CMake 3.15 or newer, make, git, Python 3 (for the Z3 build)

Z3 must be built locally. The Go bindings ship inside the Z3 source tree at
`src/api/go` and are not published to the module proxy.

## Build Z3

```sh
git clone https://github.com/Z3Prover/z3.git ../z3
cd ../z3
mkdir build && cd build
cmake -DZ3_BUILD_GO_BINDINGS=ON -DZ3_BUILD_LIBZ3_SHARED=ON ..
make -j"$(nproc)"
