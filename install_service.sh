#!/bin/bash

cat <<EOF > /etc/systemd/system/ligolo-mp.service
[Unit]
Description=Ligolo-mp
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=on-failure
RestartSec=3
User=root
ExecStart=/usr/local/bin/ligolo-mp -daemon

[Install]
WantedBy=multi-user.target
EOF

chown root:root /etc/systemd/system/ligolo-mp.service
chmod 600 /etc/systemd/system/ligolo-mp.service
systemctl daemon-reload
