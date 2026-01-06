#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeComposite

[Service]
ExecStart=/usr/local/bin/hlc_amd64_linux --path ~/.config/hidden-lake
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_composite.service

cd /usr/local/bin && \
    rm -f hlc_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hlc_amd64_linux && \
    chmod +x hlc_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_composite.service
systemctl restart hidden_lake_composite.service
