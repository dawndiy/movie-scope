#/bin/bash

rm -rf movie-scope
go build
unity-scope-tool movie-scope.ini
