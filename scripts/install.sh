#!/bin/bash



installLinux () {
    echo "Installing on Linux"
    sudo cp -r ./build/bin/reaper /usr/local/bin
    /usr/local/bin/reaper
}

installMac () {
    echo "Installing on Mac..."
    cp -r ./build/bin/reaper.app  ~/Applications
    open ~/Applications/reaper.app
}

notSupported () {
    echo "OS not supported: ${unameOut}"
    exit 1
}

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     installLinux;;
    Darwin*)    installMac;;
    *)          notSupported
esac
