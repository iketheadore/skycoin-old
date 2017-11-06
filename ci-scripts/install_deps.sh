#!/usr/bin/env bash

go get github.com/mitchellh/gox
mkdir -p $GOPATH/src/github.com/skycoin
cp -rf $GOPATH/src/github.com/iketheadore/skycoin $GOPATH/src/github.com/skycoin/

if [[ $TRAVIS_OS_NAME == 'linux' ]]; then
    sudo apt-get install python-software-properties
    sudo add-apt-repository ppa:ubuntu-wine/ppa -y
    sudo apt-get update

    curl -sL https://deb.nodesource.com/setup_7.x | sudo -E bash -
    sudo apt-get install nodejs
    sudo npm install --global electron-builder
    sudo apt-get install --no-install-recommends -y icnsutils graphicsmagick xz-utils

    # install wine
    sudo apt-get install --no-install-recommends -y wine1.8

    # install Mono
    sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 3FA7E0328081BFF6A14DA29AA6A19B38D3D831EF
    echo "deb http://download.mono-project.com/repo/debian wheezy main" | sudo tee /etc/apt/sources.list.d/mono-xamarin.list
    sudo apt-get install --no-install-recommends -y mono-devel ca-certificates-mono

    sudo apt-get install --no-install-recommends -y gcc-multilib g++-multilib
elif [[ $TRAVIS_OS_NAME == 'osx' ]]; then
    sudo npm install --global electron-builder 
fi