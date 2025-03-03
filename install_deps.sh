#!/bin/bash

if command -v apt-get &> /dev/null; then # Debian-based OS (Debian, Ubuntu, etc)
    echo "Installing dependencies using apt..."
    DEBIAN_FRONTEND=noninteractive apt-get update -yqq && apt-get install -yqq \
        gpg curl build-essential git \
        mingw-w64 binutils-mingw-w64 g++-mingw-w64
        INSTALLER=(apt-get install -yqq)
elif command -v yum &> /dev/null; then # Redhat-based OS (Fedora, CentOS, RHEL)
    echo "Installing dependencies using yum..."
    yum -y install gnupg curl gcc gcc-c++ make mingw64-gcc git
        INSTALLER=(yum -y)
elif command -v pacman &>/dev/null; then # Arch-based (Manjaro, Garuda, Blackarch)
        echo "Installing dependencies using pacman..."
        pacman --noconfirm -S mingw-w64-gcc mingw-w64-binutils mingw-w64-headers
    INSTALLER=(pacman --noconfirm -S)
else
    echo "Unsupported OS, exiting"
    exit
fi

# Verify if necessary tools are installed
for cmd in curl awk gpg; do
    if ! command -v "$cmd" &> /dev/null; then
        echo "$cmd could not be found, installing..."
                ${INSTALLER[@]} "$cmd"
    fi
done
