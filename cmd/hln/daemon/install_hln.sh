#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeNotifier

[Service]
ExecStart=/root/hln_amd64_linux --path /root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_notifier.service

cd /root && \
    rm -f hln_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hln_amd64_linux && \
    chmod +x hln_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_notifier.service
systemctl restart hidden_lake_notifier.service
