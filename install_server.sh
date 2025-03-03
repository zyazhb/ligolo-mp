#!/bin/bash

if [[ "$EUID" -ne 0 ]]; then
    echo "Please run as root"
    exit
fi

echo "Installing dependencies..."
./install_deps.sh

ARCH=$(uname -m)
case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="arm";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
esac

OLDCWD=$(pwd)
cd /root || exit
echo "Running from $(pwd)"

echo "Fetching latest ligolo-mp release..."
ARTIFACTS=$(curl -s "https://api.github.com/repos/ttpreport/ligolo-mp/releases/latest" | awk -F '"' '/browser_download_url/{print $4}')
SERVER_BINARY="ligolo-mp_linux_${ARCH}"
CHECKSUMS_FILE="ligolo-mp_checksums.txt"


for URL in $ARTIFACTS
do
    if [[ "$URL" == *"$SERVER_BINARY"* ]]; then
        echo "Downloading $URL"
        curl --silent -L "$URL" --output "$(basename "$URL")"
    fi
    if [[ "$URL" == *"$CHECKSUMS_FILE"* ]]; then
        echo "Downloading $URL"
        curl --silent -L "$URL" --output "$(basename "$URL")"
    fi
done

# Signature verification
echo "Verifying signatures ..."
sha256sum --ignore-missing -c ligolo-mp_checksums.txt || (echo "Signature mismatch! Aborting..." && exit 2)
echo

if test -f "/root/$SERVER_BINARY"; then
    echo "Moving the executable to /usr/local/bin/ligolo-mp..."
    mv "/root/$SERVER_BINARY" /usr/local/bin/ligolo-mp

    echo "Setting permissions for the server executable..."
    chmod 755 /usr/local/bin/ligolo-mp
    echo
else
    echo "$SERVER_BINARY not found! Aborting..." 
    exit 3
fi

echo "Configuring systemd service ..."
cd $OLDCWD
./install_service.sh
echo

echo "Starting the Ligolo-mp service..."
systemctl enable ligolo-mp
systemctl restart ligolo-mp
echo
