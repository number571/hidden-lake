#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakePinger

[Service]
ExecStart=/root/hls-pinger_amd64_linux --path /root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_service_pinger.service

cd /root && \
    rm -f hls-pinger_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hls-pinger_amd64_linux && \
    chmod +x hls-pinger_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_service_pinger.service
systemctl restart hidden_lake_service_pinger.service
