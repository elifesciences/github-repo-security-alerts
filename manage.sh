#!/bin/bash
set -e

cmd="$1"

if test ! "$cmd"; then
    echo "command required."
    echo
    echo "available commands:"
    echo "  build        build project for Linux"
    echo "  update-deps  update project dependencies"
    exit 1
fi

shift
rest=$*

if test "$cmd" = "build"; then
    # GOOS is 'Go OS' and is being explicit in which OS to build for
    # ld -s is 'disable symbol table'
    # ld -w is 'disable DWARF generation'
    # -v 'verbose'
    # -race 'data race detection'
    GOOS=linux go build -ldflags="-s -w" \
        -v \
        -race
    exit 0

elif test "$cmd" = "update-deps"; then
    # -u 'update modules [...] to use newer minor or patch releases when available'
    go get -u
    go mod tidy
    exit 0

# ...

fi

echo "unknown command: $cmd"
exit 1
