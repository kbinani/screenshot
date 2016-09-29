#!/bin/bash

(
	cd "$HOME"
	if [ ! -d go ]; then
		wget https://storage.googleapis.com/golang/go1.7.1.darwin-amd64.tar.gz
		tar zxf go1.7.1.darwin-amd64.tar.gz
	fi
	if [ ! -d gobootstrap ]; then
		cp -R go gobootstrap
	fi
)
