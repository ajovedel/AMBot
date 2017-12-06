#!/bin/sh

git pull
go get -u -v all
go -v build
./AMBot
