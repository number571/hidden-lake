#!/bin/bash

echo "
[Unit]
Description=HiddenLakeAdapterHTTP

[Service]
ExecStart=$HOME/.hidden-lake/bin/hla-http_amd64_linux --path $HOME/hidden-lake/etc
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
" > ~/.config/systemd/user/hidden_lake_adapter_http.service

mkdir -p ~/.hidden-lake/bin
cd ~/.hidden-lake/bin && \
    rm -f hla-http_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hla-http_amd64_linux && \
    chmod +x hla-http_amd64_linux

systemctl --user daemon-reload
systemctl --user enable hidden_lake_adapter_http.service
systemctl --user restart hidden_lake_adapter_http.service
