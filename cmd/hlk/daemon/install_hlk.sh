#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeService

[Service]
ExecStart=/usr/local/bin/hls_amd64_linux --path ~/.config/hidden-lake
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_kernel.service

cd /usr/local/bin && \
    rm -f hls_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hlk_amd64_linux && \
    chmod +x hls_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_kernel.service
systemctl restart hidden_lake_kernel.service
