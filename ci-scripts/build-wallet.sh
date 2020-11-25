#!/usr/bin/env bash

set -e -ox pipefail


if [[ "$OS_NAME" == "macOS" ]] && [[ ! "${SIGN_BRANCHES[@]}" =~ "${BRANCH}" || "$GITHUB_EVENT_NAME" == "pull_request" ]]; then
    export CSC_IDENTITY_AUTO_DISCOVERY=false;
fi

echo "start to build wallets..."
pushd "electron" >/dev/null
#echo "os name: $OS_NAME"
#if [[ "$OS_NAME" == "Linux" ]]; then ./build.sh 'linux/amd64 linux/arm' ;fi
if [[ "$OS_NAME" == "macOS" ]]; then ./build.sh ;fi
#if [[ "$OS_NAME" == "macOS" ]]; then ./build.sh 'darwin/amd64' ;fi
#if [[ "$OS_NAME" == "Windows" ]];  then ./build.sh 'windows/amd64 windows/386'; fi
ls release/
popd >/dev/null
