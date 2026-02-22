#!/bin/bash

echo "
[Unit]
Description=HiddenLakeComposite

[Service]
ExecStart=$HOME/.hidden-lake/bin/hlc_amd64_linux --path $HOME/.hidden-lake/etc
Restart=always
RestartSec=10

[Install]
WantedBy=default.target
" > ~/.config/systemd/user/hidden_lake_composite.service

mkdir -p ~/.hidden-lake/bin
cd ~/.hidden-lake/bin && \
    rm -f hlc_amd64_linux && \
    wget https://github.com/number571/hidden-lake/releases/latest/download/hlc_amd64_linux && \
    chmod +x hlc_amd64_linux

systemctl --user daemon-reload
systemctl --user enable hidden_lake_composite.service
systemctl --user restart hidden_lake_composite.service
