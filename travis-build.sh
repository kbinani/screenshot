#!/bin/bash

[ -n "$TEMP" ] || export TEMP=/tmp
export GOROOT="$HOME/go"
export GOROOT_BOOTSTRAP="$HOME/gobootstrap"
export GOPATH="$TEMP/gopath"
export PATH="$GOROOT/bin:$PATH"

env | grep -e ^GO -e ^PATH -e ^TRAVIS -e ^TEMP | sort

[ -n "$TRAVIS_BUILD_DIR" ] || exit 2
[ -n "$TRAVIS_REPO_SLUG" ] || exit 3
[ -n "$GOOS" ] || exit 4
[ -n "$GOARCH" ] || exit 5

[ -d "$TEMP" ] || exit 1
[ -d "$GOROOT" ] || exit 6
[ -d "$GOROOT_BOOTSTRAP" ] || exit 7

go version || exit 8

# setup cross compile environment
(
	cd "$GOROOT/src"
	./make.bash || exit 1
) || exit 9

# copy workspace files to GOPATH
rm -rf "$GOPATH"
mkdir -p "$GOPATH/src/github.com/$TRAVIS_REPO_SLUG" || exit 10
cp -R "$TRAVIS_BUILD_DIR" "$GOPATH/src/github.com/$TRAVIS_REPO_SLUG/.." || exit 11

# install dependencies
if [ "$GOOS" = "linux" -o "$GOOS" = "freebsd" ]; then
	go get github.com/BurntSushi/xgb || exit 12
elif [ "$GOOS" = "windows" ]; then
	go get github.com/lxn/win || exit 13
fi

# build example/main.go
(
	cd "$GOPATH/src/github.com/$TRAVIS_REPO_SLUG"
	go build example/main.go || exit 1
	echo "Built successfully"
	ls -la
) || exit 14
