#/bin/bash

./build-click-package.sh movie-scope ubuntu-sdk-15.04 vivid
adb push movie-scope.ubuntu-dawndiy_0.1.1_armhf.click /home/phablet
adb shell 'pkcon install-local movie-scope.ubuntu-dawndiy_0.1.1_armhf.click --allow-untrusted'
