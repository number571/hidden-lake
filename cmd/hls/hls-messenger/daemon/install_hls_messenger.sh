#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeMessenger

[Service]
ExecStart=/usr/local/bin/hls-messenger_amd64_linux --path ~/.config/hidden-lake
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_service_messenger.service

cd /usr/local/bin && \
    rm -f hls-messenger_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hls-messenger_amd64_linux && \
    chmod +x hls-messenger_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_service_messenger.service
systemctl restart hidden_lake_service_messenger.service
