#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeAdapterTCP

[Service]
ExecStart=/root/hla-tcp_amd64_linux --path /root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_adapter_tcp.service

cd /root && \
    rm -f hla-tcp_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hla-tcp_amd64_linux && \
    chmod +x hla-tcp_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_adapter_tcp.service
systemctl restart hidden_lake_adapter_tcp.service
