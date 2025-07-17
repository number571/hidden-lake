#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeRemoter

[Service]
ExecStart=/root/hls_remoter_amd64_linux --path /root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_service_remoter.service

cd /root && \
    rm -f hls_remoter_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hls_remoter_amd64_linux && \
    chmod +x hls_remoter_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_service_remoter.service
systemctl restart hidden_lake_service_remoter.service
