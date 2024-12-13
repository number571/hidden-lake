#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeAdapterTCP

[Service]
ExecStart=/root/hla_tcp_amd64_linux -path=/root -threads=1
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_adapter_tcp.service

cd /root && \
    rm -f hla_tcp_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hla_tcp_amd64_linux && \
    chmod +x hla_tcp_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_adapter_tcp.service
systemctl restart hidden_lake_adapter_tcp.service
