#!/bin/bash

[ -n "$GO_VERSION" ] || exit 1

(
	cd "$HOME"
	if [ ! -f "go${GO_VERSION}/bin/go" ]; then
		rm -rf "go${GO_VERSION}"
		rm -rf "gobootstrap${GO_VERSION}"

		wget "https://storage.googleapis.com/golang/go${GO_VERSION}.darwin-amd64.tar.gz" -O "go${GO_VERSION}.darwin-amd64.tar.gz" || exit 2
		rm -rf go
		tar zxf "go${GO_VERSION}.darwin-amd64.tar.gz" || exit 3

		mv go "go${GO_VERSION}" || exit 4
		rm "go${GO_VERSION}.darwin-amd64.tar.gz"
		cp -R "go${GO_VERSION}" "gobootstrap${GO_VERSION}" || exit 5
	fi
) || exit $?
