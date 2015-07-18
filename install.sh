#/bin/bash

./build-click-package.sh movie-scope ubuntu-sdk-15.04 vivid
adb push movie-scope_0.1_armhf.click /home/phablet
adb shell 'pkcon install-local movie-scope_0.1_armhf.click --allow-untrusted'
