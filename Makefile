build:
	CGO_CFLAGS="-I../z3/src/api" \
	CGO_LDFLAGS="../z3/build/libz3.so -lpthread" \
	go build -ldflags="-r ../z3/build" -o bin/z3 cmd/z3/main.go

clean:
	rm -rf bin
