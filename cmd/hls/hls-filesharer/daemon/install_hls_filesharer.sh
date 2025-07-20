#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeFilesharer

[Service]
ExecStart=/root/hls-filesharer_amd64_linux --path /root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_service_filesharer.service

cd /root && \
    rm -f hls-filesharer_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hls-filesharer_amd64_linux && \
    chmod +x hls-filesharer_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_service_filesharer.service
systemctl restart hidden_lake_service_filesharer.service
