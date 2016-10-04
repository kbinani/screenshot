#!/bin/bash

[ -n "$GO_VERSION" ] || exit 1
[ -n "$TRAVIS_REPO_SLUG" ] || exit 6

(
	cd "$HOME"
	mkdir -p "cache/$TRAVIS_REPO_SLUG"
	if [ ! -f "cache/${TRAVIS_REPO_SLUG}/go${GO_VERSION}.darwin-amd64.tar.gz" ]; then
		wget "https://storage.googleapis.com/golang/go${GO_VERSION}.darwin-amd64.tar.gz" -O "cache/${TRAVIS_REPO_SLUG}/go${GO_VERSION}.darwin-amd64.tar.gz" || exit 2
	fi

	if [ ! -f "go${GO_VERSION}/bin/go" ]; then
		rm -rf "go${GO_VERSION}"
		rm -rf "gobootstrap${GO_VERSION}"

		rm -rf go
		tar zxf "cache/${TRAVIS_REPO_SLUG}/go${GO_VERSION}.darwin-amd64.tar.gz" || exit 3

		mv go "go${GO_VERSION}" || exit 4
		cp -R "go${GO_VERSION}" "gobootstrap${GO_VERSION}" || exit 5
	fi
) || exit $?
