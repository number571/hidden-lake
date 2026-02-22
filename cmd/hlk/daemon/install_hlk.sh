#!/bin/bash

echo "
[Unit]
Description=HiddenLakeKernel

[Service]
ExecStart=$HOME/.hidden-lake/bin/hls_amd64_linux --path $HOME/hidden-lake/etc
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
" > ~/.config/systemd/user/hidden_lake_kernel.service

mkdir -p ~/.hidden-lake/bin
cd ~/.hidden-lake/bin && \
    rm -f hls_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hlk_amd64_linux && \
    chmod +x hls_amd64_linux

systemctl --user daemon-reload
systemctl --user enable hidden_lake_kernel.service
systemctl --user restart hidden_lake_kernel.service
