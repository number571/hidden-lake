#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeAdapterHTTP

[Service]
ExecStart=/usr/local/bin/hla-http_amd64_linux --path ~/.config/hidden-lake
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_adapter_http.service

cd /usr/local/bin && \
    rm -f hla-http_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hla-http_amd64_linux && \
    chmod +x hla-http_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_adapter_http.service
systemctl restart hidden_lake_adapter_http.service
